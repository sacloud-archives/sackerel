package sacloud

import (
	"github.com/sacloud/sackerel/job/core"
	"github.com/yamamoto-febc/libsacloud/api"
	"reflect"
)

// DetectLoadBalancerJob さくらのクラウド上からロードバランサーアプライアンスを検出するジョブ
func DetectLoadBalancerJob(payload interface{}) core.JobAPI {
	return core.NewJob("DetectLoadBalancer", detectLoadBalancer, payload)
}

func detectLoadBalancer(queue *core.Queue, option *core.Option, job core.JobAPI) {

	err := doActionPerZone(option, func(client *api.Client) error {

		res, err := client.LoadBalancer.Find()
		if err != nil {
			return err
		}

		for _, loadbalancer := range res.LoadBalancers {

			if hasIgnoreTagWithInfoLog(loadbalancer, option.IgnoreTag, queue, loadbalancer.ID, loadbalancer.Name, reflect.TypeOf(loadbalancer)) {
				continue
			}

			s := loadbalancer
			payload := core.NewCreateHostPayload(&s, client.Zone, loadbalancer.ID, reflect.TypeOf(loadbalancer))

			if loadbalancer.IsAvailable() {
				if loadbalancer.Instance != nil && loadbalancer.Instance.IsUp() {
					payload.MackerelHostStatus = core.MackerelHostStatusWorking
				} else {
					payload.MackerelHostStatus = core.MackerelHostStatusPowerOff
				}
			} else {
				payload.MackerelHostStatus = core.MackerelHostStatusMaintenance
			}

			queue.PushRequest("found-loadbalancer", payload)
		}

		return nil
	})

	if err != nil {
		queue.PushError(err)
		return
	}
}
