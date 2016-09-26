package sacloud

import (
	"github.com/sacloud/sackerel/job/core"
	"github.com/yamamoto-febc/libsacloud/api"
	"reflect"
)

// DetectRouterJob さくらのクラウド上からルーターを検出するジョブ
func DetectRouterJob(payload interface{}) core.JobAPI {
	return core.NewJob("DetectRouter", detectRouter, payload)
}

func detectRouter(queue *core.Queue, option *core.Option, job core.JobAPI) {

	err := doActionPerZone(option, func(client *api.Client) error {

		res, err := client.Internet.Find()
		if err != nil {
			return err
		}

		for _, router := range res.Internet {

			// router.Switchでは情報が不足しているため、Switchのみ再取得する
			sw, err := client.Switch.Read(router.Switch.ID)
			if err != nil {
				return err
			}
			router.Switch = sw

			// !!HACK!! routerではなく、スイッチがタグを持っている
			if hasIgnoreTagWithInfoLog(sw, option.IgnoreTag, queue, router.ID, router.Name, reflect.TypeOf(router)) {
				continue
			}

			r := router
			payload := core.NewCreateHostPayload(&r, client.Zone, router.ID, reflect.TypeOf(router))
			payload.MackerelHostStatus = core.MackerelHostStatusWorking

			queue.PushRequest("found-router", payload)
		}

		return nil
	})

	if err != nil {
		queue.PushError(err)
		return
	}
}
