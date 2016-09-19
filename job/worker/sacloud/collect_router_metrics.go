package sacloud

import (
	"fmt"
	"github.com/sacloud/sackerel/job/core"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"time"
)

// CollectRouterMetricsAllJob 過去分含めたルーターメトリクス(Interface)を取得するジョブ
func CollectRouterMetricsAllJob(payload interface{}) core.JobAPI {
	return core.NewJob("CollectMetricsRouterAll", collectRouterMetricsAll, payload)
}

// CollectRouterMetricsLatestJob 直近のルーターメトリクス(Interface)を取得するジョブ
func CollectRouterMetricsLatestJob(payload interface{}) core.JobAPI {
	return core.NewJob("CollectMetricsRouterLatest", collectRouterMetricsLatest, payload)
}

func collectRouterMetricsAll(queue *core.Queue, option *core.Option, job core.JobAPI) {
	collectRouterMetricsInner(queue, option, job, nil)
}

func collectRouterMetricsLatest(queue *core.Queue, option *core.Option, job core.JobAPI) {

	end := time.Now()
	start := end.Add(-option.MetricsHistoryPeriod)
	req := sacloud.NewResourceMonitorRequest(&start, &end)
	collectRouterMetricsInner(queue, option, job, req)
}

func collectRouterMetricsInner(queue *core.Queue, option *core.Option, job core.JobAPI, req *sacloud.ResourceMonitorRequest) {
	var payload = job.GetPayload()
	if payload == nil {
		queue.PushWarn(fmt.Errorf("'%s' => payload is nil", job.GetName()))
		return
	}

	// check payload types
	var sourcePayload *core.SourcePayload

	if s, ok := payload.(core.SourcePayloadHolder); ok {
		sourcePayload = s.GetSourcePayload()
		source := sourcePayload.SacloudSource
		if _, ok := source.(*sacloud.Internet); !ok {
			queue.PushWarn(fmt.Errorf("'%s' => payload.Source is invalid type. need [*sacloud.Router]", job.GetName()))
			return
		}
	} else {
		queue.PushWarn(fmt.Errorf("'%s' => payload is invalid type. need [core.SourcePayloadHolder]", job.GetName()))
		return
	}

	// Create the collect-metrics payload
	metricsPayload := core.NewCollectMetricsPayload(sourcePayload)
	client := getClient(option, sourcePayload.SacloudZone)

	// interfaces
	nicMetrics, err := client.Internet.Monitor(sourcePayload.SacloudResourceID, req)
	if err != nil {
		queue.PushError(err)
		return
	}
	metricsPayload.Metrics.Interface = append(metricsPayload.Metrics.Interface, nicMetrics)

	queue.PushRequest("collected-router-metrics", metricsPayload)
}
