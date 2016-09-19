package sacloud

import (
	"github.com/sacloud/sackerel/job/core"
	"github.com/yamamoto-febc/libsacloud/api"
	"reflect"
)

// DetectServerJob さくらのクラウド上からサーバーを検出するジョブ
func DetectServerJob(payload interface{}) core.JobAPI {
	return core.NewJob("DetectServer", detectServer, payload)
}

func detectServer(queue *core.Queue, option *core.Option, job core.JobAPI) {

	err := doActionPerZone(option, func(client *api.Client) error {

		res, err := client.Server.Find()
		if err != nil {
			return err
		}

		for _, server := range res.Servers {

			if hasIgnoreTagWithInfoLog(server, option.IgnoreTag, queue, server.ID, server.Name, reflect.TypeOf(server)) {
				continue
			}

			s := server
			payload := core.NewCreateHostPayload(&s, client.Zone, server.ID, reflect.TypeOf(server))
			if server.IsAvailable() {
				if server.Instance != nil && server.Instance.IsUp() {
					payload.MackerelHostStatus = core.MackerelHostStatusWorking
				} else {
					payload.MackerelHostStatus = core.MackerelHostStatusPowerOff
				}
			} else {
				payload.MackerelHostStatus = core.MackerelHostStatusMaintenance
			}

			queue.PushRequest("found-server", payload)
		}

		return nil
	})

	if err != nil {
		queue.PushError(err)
		return
	}
}
