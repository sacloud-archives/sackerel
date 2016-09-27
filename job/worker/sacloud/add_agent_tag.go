package sacloud

import (
	"fmt"
	"github.com/sacloud/sackerel/job/core"
)

// AddAgentTagJob エージェントタグ追加ジョブ
func AddAgentTagJob(payload interface{}) core.JobAPI {
	return core.NewJob("AddAgentTag", addAgentTag, payload)
}

func addAgentTag(queue *core.Queue, option *core.Option, job core.JobAPI) {

	var payload = job.GetPayload()
	if payload == nil {
		queue.PushWarn(fmt.Errorf("'%s' => payload is nil", job.GetName()))
		return
	}

	// check payload types
	if reconcilePayload, ok := payload.(*core.ReconcileHostsPayload); ok {

		zone, targetID, err := reconcilePayload.GetSacloudServerInfo()
		if err != nil {
			queue.PushError(err)
			return
		}

		// sacloudから検索
		client := getClient(option, zone)
		server, err := client.Server.Read(targetID)
		if err != nil {
			queue.PushError(err)
			return
		}

		if !server.HasTag(option.AgentTag) {
			server.AppendTag(option.AgentTag)
		}
		_, err = client.Server.Update(targetID, server)
		if err != nil {
			queue.PushError(err)
			return
		}

		queue.PushRequest("added-agent-tag", payload)

	} else {
		queue.PushWarn(fmt.Errorf("'%s' => payload is invalid type. need [*core.ReconcileHostsPayload]", job.GetName()))
		return
	}

}
