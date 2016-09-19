package mackerel

import (
	"fmt"
	mkr "github.com/mackerelio/mackerel-client-go"
	"github.com/sacloud/sackerel/job/core"
)

func getClient(option *core.Option) *mkr.Client {
	if option.MackerelOption.TraceMode {
		c, e := mkr.NewClientWithOptions(option.MackerelOption.APIKey, "https://mackerel.io/api/v0", true)
		if e != nil {
			panic(e)
		}
		c.UserAgent = "sackerel-trace-mode"
		return c
	}

	c := mkr.NewClient(option.MackerelOption.APIKey)
	c.UserAgent = "sackerel"
	return c

}

func findMackerelTarget(queue *core.Queue, option *core.Option, job core.JobAPI, resourceTypeName string) {

	var payload = job.GetPayload()
	if payload == nil {
		queue.PushWarn(fmt.Errorf("'%s' => payload is nil", job.GetName()))
		return
	}

	if findPayload, ok := payload.(core.MackerelFindParamHolder); ok {
		findParam := findPayload.GetFindParam()
		if findParam == nil {
			queue.PushWarn(fmt.Errorf("'%s' => payload is need FindParam", job.GetName()))
			return
		}

		// find from mackerel
		client := getClient(option)
		hosts, err := client.FindHosts(findParam)
		if err != nil {
			queue.PushError(err)
			return
		}

		if len(hosts) == 0 {
			queue.PushRequest("mackerel-not-found-"+resourceTypeName, payload)
		} else {
			for _, host := range hosts {

				//IDの設定
				if source, ok := payload.(core.SourcePayloadHolder); ok {
					mackerelInfo := source.GetSourcePayload()
					mackerelInfo.MackerelID = host.ID
					mackerelInfo.MackerelHost = host
				}

				queue.PushRequest("mackerel-found-"+resourceTypeName, payload)
			}
		}

	} else {

		queue.PushWarn(fmt.Errorf("'%s' => payload is invalid type. need [core.MackerelFindParamHolder]", job.GetName()))
		return
	}

}

func registMackerelTarget(queue *core.Queue, option *core.Option, job core.JobAPI, resourceTypeName string) {

	var payload = job.GetPayload()
	if payload == nil {
		queue.PushWarn(fmt.Errorf("'%s' => payload is nil", job.GetName()))
		return
	}

	if CreateHostPayload, ok := payload.(*core.CreateHostPayload); ok {
		createParam := CreateHostPayload.MackerelHostParam
		if createParam == nil {
			queue.PushWarn(fmt.Errorf("'%s' => payload is need CreateHostParam", job.GetName()))
			return
		}

		// create host in mackerel

		client := getClient(option)
		id, err := client.CreateHost(createParam)

		if err != nil {
			queue.PushError(err)
			return
		}

		targetName := fmt.Sprintf("(%s)", CreateHostPayload.SourceType.Name())
		CreateHostPayload.MackerelID = id

		host, err := client.FindHost(id)
		if err != nil {
			queue.PushError(err)
		}
		if host == nil {
			queue.PushError(fmt.Errorf("Created host not found!"))
			return
		}
		CreateHostPayload.MackerelHost = host

		queue.PushInfo(
			fmt.Sprintf(
				"Host created   %-15s => SakuraID:[%d] / MackerelID:[%s] / Name:[%s]",
				targetName,
				CreateHostPayload.SacloudResourceID,
				id,
				createParam.DisplayName,
			),
		)
		queue.PushRequest("mackerel-registed-"+resourceTypeName, payload)

	} else {

		queue.PushWarn(fmt.Errorf("'%s' => payload is invalid type. need [core.CreateHostPayload]", job.GetName()))
		return
	}

}
