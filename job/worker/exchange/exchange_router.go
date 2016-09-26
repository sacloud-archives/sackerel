package exchange

import (
	"fmt"
	mkr "github.com/mackerelio/mackerel-client-go"
	"github.com/sacloud/sackerel/job/core"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"net"
	"strings"
)

// RouterJob ルーター変換用ジョブ
func RouterJob(payload interface{}) core.JobAPI {
	return core.NewJob("ExchangeRouter", exchangeRouter, payload)
}

func exchangeRouter(queue *core.Queue, option *core.Option, job core.JobAPI) {

	var payload = job.GetPayload()
	if payload == nil {
		queue.PushWarn(fmt.Errorf("'%s' => payload is nil", job.GetName()))
		return
	}

	if serverPayload, ok := payload.(*core.CreateHostPayload); ok {
		if router, ok := serverPayload.SacloudSource.(*sacloud.Internet); ok {
			// exchange [sacloud.Server => mkr.CreateHostParam]
			hostParam, err := exchangeRouterToMackerelHost(
				serverPayload.GenerateMackerelName(),
				serverPayload.SacloudZone,
				router)

			if err != nil {
				queue.PushError(err)
				return
			}

			serverPayload.MackerelHostParam = hostParam
			// trigger next route
			queue.PushRequest("exchanged-router", serverPayload)

		} else {
			queue.PushWarn(fmt.Errorf("'%s' => payload.Source is invalid type. need [*sacloud.Server]", job.GetName()))
			return
		}

	} else {
		queue.PushWarn(fmt.Errorf("'%s' => payload is invalid type. need [*core.CreateHostPayload]", job.GetName()))
		return
	}

}

func exchangeRouterToMackerelHost(mackerelName string, zone string, router *sacloud.Internet) (*mkr.CreateHostParam, error) {

	p := &mkr.CreateHostParam{
		Name:             mackerelName,
		DisplayName:      router.Name,
		CustomIdentifier: mackerelName,
		RoleFullnames: []string{
			"SakuraCloud:Router",
			fmt.Sprintf("SakuraCloud:Zone-%s", zone),
		},
	}

	// !!HACK!!
	for _, tag := range router.Switch.Tags {
		if tag != "" && !strings.HasPrefix(tag, "@") {
			p.RoleFullnames = append(p.RoleFullnames, fmt.Sprintf("SakuraCloud:%s", tag))
		}
	}

	ipaddresses, err := getIPAddressListFromSwitchSubnet(router.Switch)
	if ipaddresses == nil || err != nil {
		return nil, err
	}
	for i, ip := range ipaddresses {
		p.Interfaces = append(p.Interfaces, mkr.Interface{
			Name:       fmt.Sprintf("eth%d", i),
			IPAddress:  ip,
			MacAddress: "",
		})
	}

	return p, nil
}

func getIPAddressListFromSwitchSubnet(s *sacloud.Switch) ([]string, error) {
	if s.Subnets == nil || len(s.Subnets) < 1 {
		return nil, fmt.Errorf("switch[%s].Subnets is nil", s.ID)
	}

	//さくらのクラウドの仕様上/24までしか割り当てできないためこのロジックでOK
	baseIP := net.ParseIP(s.Subnets[0].IPAddresses.Min).To4()
	min := baseIP[3]
	max := net.ParseIP(s.Subnets[0].IPAddresses.Max).To4()[3]

	var i byte
	ret := []string{}
	for (min + i) <= max { //境界含む
		ip := net.IPv4(baseIP[0], baseIP[1], baseIP[2], baseIP[3]+i)
		ret = append(ret, ip.String())
		i++
	}

	return ret, nil
}
