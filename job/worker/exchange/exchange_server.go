package exchange

import (
	"fmt"
	mkr "github.com/mackerelio/mackerel-client-go"
	"github.com/sacloud/libsacloud/sacloud"
	"github.com/sacloud/sackerel/job/core"
	"strings"
)

// ServerJob サーバー変換用ジョブ
func ServerJob(payload interface{}) core.JobAPI {
	return core.NewJob("ExchangeServer", exchangeServer, payload)
}

func exchangeServer(queue *core.Queue, option *core.Option, job core.JobAPI) {

	var payload = job.GetPayload()
	if payload == nil {
		queue.PushWarn(fmt.Errorf("'%s' => payload is nil", job.GetName()))
		return
	}

	if serverPayload, ok := payload.(*core.CreateHostPayload); ok {
		if server, ok := serverPayload.SacloudSource.(*sacloud.Server); ok {
			// exchange [sacloud.Server => mkr.CreateHostParam]
			serverPayload.MackerelHostParam = exchangeServerToMackerelHost(
				serverPayload.GenerateMackerelName(),
				serverPayload.SacloudZone,
				server)

			// trigger next route
			queue.PushRequest("exchanged-server", serverPayload)

		} else {
			queue.PushWarn(fmt.Errorf("'%s' => payload.Source is invalid type. need [*sacloud.Server]", job.GetName()))
			return
		}

	} else {
		queue.PushWarn(fmt.Errorf("'%s' => payload is invalid type. need [*core.CreateHostPayload]", job.GetName()))
		return
	}

}

func exchangeServerToMackerelHost(mackerelName string, zone string, server *sacloud.Server) *mkr.CreateHostParam {

	p := &mkr.CreateHostParam{
		Name:             mackerelName,
		DisplayName:      server.Name,
		CustomIdentifier: mackerelName,
		RoleFullnames: []string{
			"SakuraCloud:Server",
			fmt.Sprintf("SakuraCloud:Zone-%s", zone),
		},
	}

	for _, tag := range server.Tags {
		if tag != "" && !strings.HasPrefix(tag, "@") {
			p.RoleFullnames = append(p.RoleFullnames, fmt.Sprintf("SakuraCloud:%s", tag))
		}
	}

	for i, nic := range server.Interfaces {

		ip := ""
		if nic.Switch != nil {
			ip = nic.IPAddress
			if nic.Switch.Scope != sacloud.ESCopeShared {
				ip = nic.UserIPAddress
			}
		}

		p.Interfaces = append(p.Interfaces, mkr.Interface{
			Name:       fmt.Sprintf("eth%d", i),
			IPAddress:  ip,
			MacAddress: nic.MACAddress,
		})
	}

	return p
}
