package mackerel

import (
	mkr "github.com/mackerelio/mackerel-client-go"
	"github.com/sacloud/sackerel/job/core"
	"strings"
)

// ReconcileByMACAddressJob Mackerel上のホストから同一のMACアドレスを持つホストを検索
func ReconcileByMACAddressJob(payload interface{}) core.JobAPI {
	return core.NewJob("MackerelReconcileHosts", reconcileByMACAddress, payload)
}

func reconcileByMACAddress(queue *core.Queue, option *core.Option, job core.JobAPI) {

	// 全ホスト取得。重いかもしれない REMARKS[要改善]
	client := getClient(option)
	findOption := &mkr.FindHostsParam{
		Statuses: []string{
			string(core.MackerelHostStatusWorking),
			string(core.MackerelHostStatusStandby),
			string(core.MackerelHostStatusPowerOff),
			string(core.MackerelHostStatusMaintenance),
		},
	}
	allHosts, err := client.FindHosts(findOption)

	if err != nil {
		queue.PushError(err)
		return
	}

	for _, target := range allHosts {
		if !isNeedReconcile(target) {
			continue
		}

		found, fromAgentHost := findMatchMACAddressHost(allHosts, target)
		if found {
			payload := &core.ReconcileHostsPayload{
				FromAgentHost:    fromAgentHost,
				FromSackerelHost: target,
			}
			queue.PushRequest("reconcile-host", payload)
		}
	}
}

func isNeedReconcile(target *mkr.Host) bool {
	//   - NICが1つ以上あり、MACアドレスが設定されていること
	if len(target.Interfaces) == 0 || target.Interfaces[0].MacAddress == "" {
		return false
	}

	//   - FullRollName に "SakuraCloud:Server"が存在すること
	res := false
	roleFullNames := target.GetRoleFullnames()
	for _, name := range roleFullNames {
		if name == "SakuraCloud:Server" {
			res = true
			break
		}
	}

	return res
}

// findMatchMACAddressHost 同一のMACアドレスを持つホストを検索、ヒットした最初のホストを返す
func findMatchMACAddressHost(sources []*mkr.Host, target *mkr.Host) (bool, *mkr.Host) {
	targetMACAddress := getMACAddress(target)
	for _, source := range sources {
		// 同一IDは除く
		if source.ID == target.ID {
			continue
		}
		sourceMACAddress := getMACAddress(source)
		if strings.ToLower(targetMACAddress) == strings.ToLower(sourceMACAddress) {
			return true, source
		}
	}

	return false, nil
}

func getMACAddress(target *mkr.Host) string {
	if len(target.Interfaces) == 0 {
		return ""
	}
	return target.Interfaces[0].MacAddress
}
