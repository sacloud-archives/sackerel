package mackerel

import (
	"fmt"
	mkr "github.com/mackerelio/mackerel-client-go"
	"github.com/sacloud/sackerel/job/core"
)

// UpdateHostJob Mackerelホスト情報の更新ジョブ
func UpdateHostJob(payload interface{}) core.JobAPI {
	return core.NewJob("MackerelUpdateHost", updateHost, payload)
}

// UpdateStatusJob Mackerelホストのステータス更新ジョブ
func UpdateStatusJob(payload interface{}) core.JobAPI {
	return core.NewJob("MackerelUpdateStatus", updateStatus, payload)
}

// UpdateRoleJob Mackerelホストのロール更新用ジョブ
func UpdateRoleJob(payload interface{}) core.JobAPI {
	return core.NewJob("MackerelUpdateRole", updateRole, payload)
}

// UpdadteCustomIDJob MackerelホストのCustomIdentifier設定用ジョブ
func UpdadteCustomIDJob(payload interface{}) core.JobAPI {
	return core.NewJob("MackerelUpdateCustomID", updateCustomID, payload)
}

// RetireJob Mackerelホストの退役処理用ジョブ
func RetireJob(payload interface{}) core.JobAPI {
	return core.NewJob("MackerelRetire", retireHost, payload)
}

func updateStatus(queue *core.Queue, option *core.Option, job core.JobAPI) {

	var payload = job.GetPayload()
	if payload == nil {
		queue.PushWarn(fmt.Errorf("'%s' => payload is nil", job.GetName()))
		return
	}

	if sourcePayload, ok := payload.(*core.CreateHostPayload); ok {

		if hasSakuraCloudTag(sourcePayload.SacloudSource, option.AgentTag) {
			return
		}

		sourceInfo := sourcePayload.GetSourcePayload()

		updated, err := sourcePayload.IsStatusUpdated()
		if err != nil {
			queue.PushError(err)
			return
		}

		if sourceInfo.MackerelID == "" || !updated {
			return
		}

		client := getClient(option)

		err = client.UpdateHostStatus(sourceInfo.MackerelID, string(sourceInfo.MackerelHostStatus))

		if err != nil {
			queue.PushError(err)
			return
		}

		targetName := fmt.Sprintf("(%s)", sourceInfo.SourceType.Name())
		queue.PushInfo(
			fmt.Sprintf("Status updated %-15s => SakuraID:[%d] / MackerelID:[%s] / Status:[%s]",
				targetName,
				sourceInfo.SacloudResourceID,
				sourceInfo.MackerelID,
				sourceInfo.MackerelHostStatus,
			),
		)

		queue.PushRequest("mackerel-updated-status", payload)
	} else {
		queue.PushWarn(fmt.Errorf("'%s' => payload is invalid type. need [*core.CreateHostPayload]", job.GetName()))
		return
	}
}

func updateHost(queue *core.Queue, option *core.Option, job core.JobAPI) {

	var payload = job.GetPayload()
	if payload == nil {
		queue.PushWarn(fmt.Errorf("'%s' => payload is nil", job.GetName()))
		return
	}

	if sourcePayload, ok := payload.(*core.CreateHostPayload); ok {

		if hasSakuraCloudTag(sourcePayload.SacloudSource, option.AgentTag) {
			return
		}

		sourceInfo := sourcePayload.GetSourcePayload()
		if sourceInfo.MackerelID != "" {
			client := getClient(option)

			updateParam := mkr.UpdateHostParam(*sourcePayload.MackerelHostParam)
			_, err := client.UpdateHost(sourceInfo.MackerelID, &updateParam)
			if err != nil {
				queue.PushError(err)
				return
			}
		}

		targetName := fmt.Sprintf("(%s)", sourceInfo.SourceType.Name())
		queue.PushInfo(
			fmt.Sprintf("Host updated %-15s   => SakuraID:[%d] / MackerelID:[%s] / Name:[%s]",
				targetName,
				sourceInfo.SacloudResourceID,
				sourceInfo.MackerelID,
				sourcePayload.MackerelHostParam.DisplayName,
			),
		)

		queue.PushRequest("mackerel-updated-host", payload)
	} else {
		queue.PushWarn(fmt.Errorf("'%s' => payload is invalid type. need [*core.CreateHostPayload]", job.GetName()))
		return
	}
}

func updateRole(queue *core.Queue, option *core.Option, job core.JobAPI) {

	var payload = job.GetPayload()
	if payload == nil {
		queue.PushWarn(fmt.Errorf("'%s' => payload is nil", job.GetName()))
		return
	}

	if sourcePayload, ok := payload.(*core.CreateHostPayload); ok {

		sourceInfo := sourcePayload.GetSourcePayload()

		updated, err := sourcePayload.IsRoleUpdated()
		if err != nil {
			queue.PushError(err)
			return
		}
		if sourceInfo.MackerelID == "" || !updated {
			return
		}

		client := getClient(option)

		err = client.UpdateHostRoleFullnames(sourceInfo.MackerelID, sourcePayload.MackerelHostParam.RoleFullnames)
		if err != nil {
			queue.PushError(err)
			return
		}

		targetName := fmt.Sprintf("(%s)", sourceInfo.SourceType.Name())
		queue.PushInfo(
			fmt.Sprintf("Role updated %-15s   => SakuraID:[%d] / MackerelID:[%s] / RoleCount:[%d]",
				targetName,
				sourceInfo.SacloudResourceID,
				sourceInfo.MackerelID,
				len(sourcePayload.MackerelHostParam.RoleFullnames),
			),
		)

		queue.PushRequest("mackerel-updated-host", payload)
	} else {
		queue.PushWarn(fmt.Errorf("'%s' => payload is invalid type. need [*core.CreateHostPayload]", job.GetName()))
		return
	}
}

func updateCustomID(queue *core.Queue, option *core.Option, job core.JobAPI) {
	var payload = job.GetPayload()
	if payload == nil {
		queue.PushWarn(fmt.Errorf("'%s' => payload is nil", job.GetName()))
		return
	}

	if sourcePayload, ok := payload.(*core.ReconcileHostsPayload); ok {

		customID := sourcePayload.FromSackerelHost.Name
		target := sourcePayload.FromAgentHost

		client := getClient(option)
		param := &mkr.UpdateHostParam{
			Name:             target.Name,
			CustomIdentifier: customID,
			DisplayName:      target.DisplayName,
			Interfaces:       target.Interfaces,
			Meta:             target.Meta,
			RoleFullnames:    target.GetRoleFullnames(),
		}

		_, err := client.UpdateHost(target.ID, param)
		if err != nil {
			queue.PushError(err)
			return
		}

		queue.PushInfo(fmt.Sprintf("CustomIdentifier is updated => MackerelID:[%s] / CustomID:[%s]", target.ID, customID))

		queue.PushRequest("mackerel-updated-custom-id", payload)
	} else {
		queue.PushWarn(fmt.Errorf("'%s' => payload is invalid type. need [*core.ReconcileHostsPayload]", job.GetName()))
		return
	}
}

func retireHost(queue *core.Queue, option *core.Option, job core.JobAPI) {
	var payload = job.GetPayload()
	if payload == nil {
		queue.PushWarn(fmt.Errorf("'%s' => payload is nil", job.GetName()))
		return
	}

	if sourcePayload, ok := payload.(*core.ReconcileHostsPayload); ok {

		target := sourcePayload.FromSackerelHost

		client := getClient(option)
		err := client.RetireHost(target.ID)
		if err != nil {
			queue.PushError(err)
			return
		}

		queue.PushInfo(fmt.Sprintf("Mackerel host is retired => MackerelID:[%s]", target.ID))

		queue.PushRequest("mackerel-retired-host", payload)
	} else {
		queue.PushWarn(fmt.Errorf("'%s' => payload is invalid type. need [*core.ReconcileHostsPayload]", job.GetName()))
		return
	}
}
