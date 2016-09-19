package sacloud

import (
	"github.com/sacloud/sackerel/job/core"
	"github.com/yamamoto-febc/libsacloud/api"
	"reflect"
)

// DetectDatabaseJob さくらのクラウド上からデータベースアプライアンスを検出するジョブ
func DetectDatabaseJob(payload interface{}) core.JobAPI {
	return core.NewJob("DetectDatabase", detectDatabase, payload)
}

func detectDatabase(queue *core.Queue, option *core.Option, job core.JobAPI) {

	err := doActionPerZone(option, func(client *api.Client) error {

		res, err := client.Database.Find()
		if err != nil {
			return err
		}

		for _, database := range res.Databases {

			if hasIgnoreTagWithInfoLog(database, option.IgnoreTag, queue, database.ID, database.Name, reflect.TypeOf(database)) {
				continue
			}

			r := database
			payload := core.NewCreateHostPayload(&r, client.Zone, database.ID, reflect.TypeOf(database))

			if database.IsAvailable() {
				if database.Instance != nil && database.Instance.IsUp() {
					payload.MackerelHostStatus = core.MackerelHostStatusWorking
				} else {
					payload.MackerelHostStatus = core.MackerelHostStatusPowerOff
				}
			} else {
				payload.MackerelHostStatus = core.MackerelHostStatusMaintenance
			}

			queue.PushRequest("found-database", payload)
		}

		return nil
	})

	if err != nil {
		queue.PushError(err)
		return
	}
}
