package exchange

import (
	"fmt"
	mkr "github.com/mackerelio/mackerel-client-go"
	"github.com/sacloud/sackerel/job/core"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"strings"
)

// DatabaseJob データベース変換用ジョブ
func DatabaseJob(payload interface{}) core.JobAPI {
	return core.NewJob("ExchangeDatabase", exchangeDatabase, payload)
}

func exchangeDatabase(queue *core.Queue, option *core.Option, job core.JobAPI) {

	var payload = job.GetPayload()
	if payload == nil {
		queue.PushWarn(fmt.Errorf("'%s' => payload is nil", job.GetName()))
		return
	}

	if databasePayload, ok := payload.(*core.CreateHostPayload); ok {
		if database, ok := databasePayload.SacloudSource.(*sacloud.Database); ok {
			// exchange [sacloud.Database => mkr.CreateHostParam]
			databasePayload.MackerelHostParam = exchangeDatabaseToMackerelHost(
				databasePayload.GenerateMackerelName(),
				databasePayload.SacloudZone,
				database)

			// trigger next route
			queue.PushRequest("exchanged-database", databasePayload)

		} else {
			queue.PushWarn(fmt.Errorf("'%s' => payload.Source is invalid type. need [*sacloud.Database]", job.GetName()))
			return
		}

	} else {
		queue.PushWarn(fmt.Errorf("'%s' => payload is invalid type. need [*core.CreateHostPayload]", job.GetName()))
		return
	}

}

func exchangeDatabaseToMackerelHost(mackerelName string, zone string, database *sacloud.Database) *mkr.CreateHostParam {

	p := &mkr.CreateHostParam{
		Name:             mackerelName,
		DisplayName:      database.Name,
		CustomIdentifier: mackerelName,
		RoleFullnames: []string{
			"SakuraCloud:Database",
			fmt.Sprintf("SakuraCloud:Zone-%s", zone),
		},
	}

	for _, tag := range database.Tags {
		if tag != "" && !strings.HasPrefix(tag, "@") {
			p.RoleFullnames = append(p.RoleFullnames, fmt.Sprintf("SakuraCloud:%s", tag))
		}
	}

	ip := database.Interfaces[0].IPAddress // 共有セグメントに接続されている場合
	if database.Interfaces[0].Switch.Scope != sacloud.ESCopeShared {
		// スイッチに接続されている場合
		ip = database.Remark.Servers[0].(map[string]interface{})["IPAddress"].(string)
	}

	p.Interfaces = append(p.Interfaces, mkr.Interface{
		Name:       fmt.Sprintf("eth%d", 0),
		IPAddress:  ip,
		MacAddress: "",
	})

	return p
}
