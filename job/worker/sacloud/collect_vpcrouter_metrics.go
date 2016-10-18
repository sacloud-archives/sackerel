package sacloud

import (
	"fmt"
	"github.com/sacloud/libsacloud/sacloud"
	"github.com/sacloud/sackerel/job/core"
	"time"
)

// CollectVPCRouterMetricsAllJob 過去分含めたVPCルーターメトリクス(Interface)を取得するジョブ
func CollectVPCRouterMetricsAllJob(payload interface{}) core.JobAPI {
	return core.NewJob("CollectMetricsVPCRouterAll", collectVPCRouterMetricsAll, payload)
}

// CollectVPCRouterMetricsLatestJob 直近のVPCルーターメトリクス(Interface)を取得するジョブ
func CollectVPCRouterMetricsLatestJob(payload interface{}) core.JobAPI {
	return core.NewJob("CollectMetricsVPCRouterLatest", collectVPCRouterMetricsLatest, payload)
}

func collectVPCRouterMetricsAll(queue *core.Queue, option *core.Option, job core.JobAPI) {
	collectVPCRouterMetricsInner(queue, option, job, nil)
}

func collectVPCRouterMetricsLatest(queue *core.Queue, option *core.Option, job core.JobAPI) {

	end := time.Now()
	start := end.Add(-option.MetricsHistoryPeriod)
	req := sacloud.NewResourceMonitorRequest(&start, &end)
	collectVPCRouterMetricsInner(queue, option, job, req)
}

func collectVPCRouterMetricsInner(queue *core.Queue, option *core.Option, job core.JobAPI, req *sacloud.ResourceMonitorRequest) {
	var payload = job.GetPayload()
	if payload == nil {
		queue.PushWarn(fmt.Errorf("'%s' => payload is nil", job.GetName()))
		return
	}

	// check payload types
	var vpcrouter *sacloud.VPCRouter
	var sourcePayload *core.SourcePayload

	if s, ok := payload.(core.SourcePayloadHolder); ok {
		sourcePayload = s.GetSourcePayload()
		source := sourcePayload.SacloudSource
		if vpcrouter, ok = source.(*sacloud.VPCRouter); !ok {
			queue.PushWarn(fmt.Errorf("'%s' => payload.Source is invalid type. need [*sacloud.VPCRouter]", job.GetName()))
			return
		}
	} else {
		queue.PushWarn(fmt.Errorf("'%s' => payload is invalid type. need [core.SourcePayloadHolder]", job.GetName()))
		return
	}

	if sourcePayload.MackerelHostStatus == core.MackerelHostStatusMaintenance {
		queue.PushWarn(fmt.Errorf("SakuraCloud resource['%d'] is still maintenance state. '%s' is skipped", sourcePayload.SacloudResourceID, job.GetName()))
		return
	}

	// Create the collect-metrics payload
	metricsPayload := core.NewCollectMetricsPayload(sourcePayload)
	client := getClient(option, sourcePayload.SacloudZone)

	// interfaces
	for i, nic := range vpcrouter.Settings.Router.Interfaces {

		if i == 0 || nic != nil { // i == 0 はグローバルNIC用、以降はスイッチと接続があれば値がある
			nicMetrics, err := client.VPCRouter.MonitorBy(sourcePayload.SacloudResourceID, i, req)
			if err != nil {
				queue.PushError(err)
				return
			}
			metricsPayload.Metrics.Interface = append(metricsPayload.Metrics.Interface, nicMetrics)
		} else {
			// 値がないとNICインデックスがずれるためnilを入れておく
			metricsPayload.Metrics.Interface = append(metricsPayload.Metrics.Interface, nil)
		}
	}

	queue.PushRequest("collected-vpcrouter-metrics", metricsPayload)
}
