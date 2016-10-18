package sacloud

import (
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/sackerel/job/core"
	"reflect"
)

// DetectVPCRouterJob さくらのクラウド上からVPCルーターアプライアンスを検出するジョブ
func DetectVPCRouterJob(payload interface{}) core.JobAPI {
	return core.NewJob("DetectVPCRouter", detectVPCRouter, payload)
}

func detectVPCRouter(queue *core.Queue, option *core.Option, job core.JobAPI) {

	err := doActionPerZone(option, func(client *api.Client) error {

		res, err := client.VPCRouter.Find()
		if err != nil {
			return err
		}

		for _, vpcrouter := range res.VPCRouters {

			if hasIgnoreTagWithInfoLog(vpcrouter, option.IgnoreTag, queue, vpcrouter.ID, vpcrouter.Name, reflect.TypeOf(vpcrouter)) {
				continue
			}

			r := vpcrouter
			payload := core.NewCreateHostPayload(&r, client.Zone, vpcrouter.ID, reflect.TypeOf(vpcrouter))
			if vpcrouter.IsAvailable() {
				if vpcrouter.Instance != nil && vpcrouter.Instance.IsUp() {
					payload.MackerelHostStatus = core.MackerelHostStatusWorking
				} else {
					payload.MackerelHostStatus = core.MackerelHostStatusPowerOff
				}
			} else {
				payload.MackerelHostStatus = core.MackerelHostStatusMaintenance
			}

			queue.PushRequest("found-vpcrouter", payload)
		}

		return nil
	})

	if err != nil {
		queue.PushError(err)
		return
	}
}
