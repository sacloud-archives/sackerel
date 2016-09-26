package sacloud

import (
	"fmt"
	"github.com/sacloud/sackerel/job/core"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"time"
)

// CollectLoadBalancerMetricsAllJob 過去分含めたルーターメトリクス(Interface)を取得するジョブ
func CollectLoadBalancerMetricsAllJob(payload interface{}) core.JobAPI {
	return core.NewJob("CollectMetricsLoadBalancerAll", collectLoadBalancerMetricsAll, payload)
}

// CollectLoadBalancerMetricsLatestJob 直近のルーターメトリクス(Interface)を取得するジョブ
func CollectLoadBalancerMetricsLatestJob(payload interface{}) core.JobAPI {
	return core.NewJob("CollectMetricsLoadBalancerLatest", collectLoadBalancerMetricsLatest, payload)
}

func collectLoadBalancerMetricsAll(queue *core.Queue, option *core.Option, job core.JobAPI) {
	collectLoadBalancerMetricsInner(queue, option, job, nil)
}

func collectLoadBalancerMetricsLatest(queue *core.Queue, option *core.Option, job core.JobAPI) {

	end := time.Now()
	start := end.Add(-option.MetricsHistoryPeriod)
	req := sacloud.NewResourceMonitorRequest(&start, &end)
	collectLoadBalancerMetricsInner(queue, option, job, req)
}

func collectLoadBalancerMetricsInner(queue *core.Queue, option *core.Option, job core.JobAPI, req *sacloud.ResourceMonitorRequest) {
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
		if _, ok := source.(*sacloud.LoadBalancer); !ok {
			queue.PushWarn(fmt.Errorf("'%s' => payload.Source is invalid type. need [*sacloud.LoadBalancer]", job.GetName()))
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
	nicMetrics, err := client.LoadBalancer.Monitor(sourcePayload.SacloudResourceID, req)
	if err != nil {
		queue.PushError(err)
		return
	}
	metricsPayload.Metrics.Interface = append(metricsPayload.Metrics.Interface, nicMetrics)

	queue.PushRequest("collected-loadbalancer-metrics", metricsPayload)
}
