package exchange

import (
	"fmt"
	mkr "github.com/mackerelio/mackerel-client-go"
	"github.com/sacloud/sackerel/job/core"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"strings"
)

// VPCRouterJob VPCルーター変換用ジョブ
func VPCRouterJob(payload interface{}) core.JobAPI {
	return core.NewJob("ExchangeVPCRouter", exchangeVPCRouter, payload)
}

func exchangeVPCRouter(queue *core.Queue, option *core.Option, job core.JobAPI) {

	var payload = job.GetPayload()
	if payload == nil {
		queue.PushWarn(fmt.Errorf("'%s' => payload is nil", job.GetName()))
		return
	}

	if vpcrouterPayload, ok := payload.(*core.CreateHostPayload); ok {
		if vpcrouter, ok := vpcrouterPayload.SacloudSource.(*sacloud.VPCRouter); ok {
			// exchange [sacloud.VPCRouter => mkr.CreateHostParam]
			vpcrouterPayload.MackerelHostParam = exchangeVPCRouterToMackerelHost(
				vpcrouterPayload.GenerateMackerelName(),
				vpcrouterPayload.SacloudZone,
				vpcrouter)

			// trigger next route
			queue.PushRequest("exchanged-vpcrouter", vpcrouterPayload)

		} else {
			queue.PushWarn(fmt.Errorf("'%s' => payload.Source is invalid type. need [*sacloud.VPCRouter]", job.GetName()))
			return
		}

	} else {
		queue.PushWarn(fmt.Errorf("'%s' => payload is invalid type. need [*core.CreateHostPayload]", job.GetName()))
		return
	}

}

func exchangeVPCRouterToMackerelHost(mackerelName string, zone string, vpcrouter *sacloud.VPCRouter) *mkr.CreateHostParam {

	p := &mkr.CreateHostParam{
		Name:             mackerelName,
		DisplayName:      vpcrouter.Name,
		CustomIdentifier: mackerelName,
		RoleFullnames: []string{
			"SakuraCloud:VPCRouter",
			fmt.Sprintf("SakuraCloud:Zone-%s", zone),
		},
	}

	for _, tag := range vpcrouter.Tags {
		if tag != "" && !strings.HasPrefix(tag, "@") {
			p.RoleFullnames = append(p.RoleFullnames, fmt.Sprintf("SakuraCloud:%s", tag))
		}
	}

	// Global IP
	ip := vpcrouter.Interfaces[0].IPAddress
	nic := vpcrouter.Settings.Router.Interfaces[0]
	if nic != nil && nic.VirtualIPAddress != "" {
		ip = nic.VirtualIPAddress
	}

	if ip != "" {
		p.Interfaces = append(p.Interfaces, mkr.Interface{
			Name:       fmt.Sprintf("eth%d", 0),
			IPAddress:  ip,
			MacAddress: "",
		})
	}

	// private ip
	for i, nic := range vpcrouter.Settings.Router.Interfaces {
		if i == 0 {
			continue
		}
		ip := nic.IPAddress[0] // 必ず存在するはず。
		if nic.VirtualIPAddress != "" {
			ip = nic.VirtualIPAddress
		}
		p.Interfaces = append(p.Interfaces, mkr.Interface{
			Name:       fmt.Sprintf("eth%d", i),
			IPAddress:  ip,
			MacAddress: "",
		})
	}

	return p
}
