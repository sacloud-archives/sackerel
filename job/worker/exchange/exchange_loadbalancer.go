package exchange

import (
	"fmt"
	mkr "github.com/mackerelio/mackerel-client-go"
	"github.com/sacloud/sackerel/job/core"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"strings"
)

// LoadBalancerJob ロードバランサー変換用ジョブ
func LoadBalancerJob(payload interface{}) core.JobAPI {
	return core.NewJob("ExchangeLoadBalancer", exchangeLoadBalancer, payload)
}

func exchangeLoadBalancer(queue *core.Queue, option *core.Option, job core.JobAPI) {

	var payload = job.GetPayload()
	if payload == nil {
		queue.PushWarn(fmt.Errorf("'%s' => payload is nil", job.GetName()))
		return
	}

	if loadbalancerPayload, ok := payload.(*core.CreateHostPayload); ok {
		if loadbalancer, ok := loadbalancerPayload.SacloudSource.(*sacloud.LoadBalancer); ok {
			// exchange [sacloud.LoadBalancer => mkr.CreateHostParam]
			loadbalancerPayload.MackerelHostParam = exchangeLoadBalancerToMackerelHost(
				loadbalancerPayload.GenerateMackerelName(),
				loadbalancerPayload.SacloudZone,
				loadbalancer)

			// trigger next route
			queue.PushRequest("exchanged-loadbalancer", loadbalancerPayload)

		} else {
			queue.PushWarn(fmt.Errorf("'%s' => payload.Source is invalid type. need [*sacloud.LoadBalancer]", job.GetName()))
			return
		}

	} else {
		queue.PushWarn(fmt.Errorf("'%s' => payload is invalid type. need [*core.CreateHostPayload]", job.GetName()))
		return
	}

}

func exchangeLoadBalancerToMackerelHost(mackerelName string, zone string, loadbalancer *sacloud.LoadBalancer) *mkr.CreateHostParam {

	p := &mkr.CreateHostParam{
		Name:             mackerelName,
		DisplayName:      loadbalancer.Name,
		CustomIdentifier: mackerelName,
		RoleFullnames: []string{
			"SakuraCloud:LoadBalancer",
			fmt.Sprintf("SakuraCloud:Zone-%s", zone),
		},
	}
	for _, tag := range loadbalancer.Tags {
		if tag != "" && !strings.HasPrefix(tag, "@") {
			p.RoleFullnames = append(p.RoleFullnames, fmt.Sprintf("SakuraCloud:%s", tag))
		}
	}

	p.Interfaces = append(p.Interfaces, mkr.Interface{
		Name:       fmt.Sprintf("eth%d", 0),
		IPAddress:  loadbalancer.Remark.Servers[0].(map[string]interface{})["IPAddress"].(string),
		MacAddress: "",
	})
	if len(loadbalancer.Remark.Servers) > 1 {
		p.Interfaces = append(p.Interfaces, mkr.Interface{
			Name:       fmt.Sprintf("eth%d", 1),
			IPAddress:  loadbalancer.Remark.Servers[1].(map[string]interface{})["IPAddress"].(string),
			MacAddress: "",
		})
	}

	return p
}
