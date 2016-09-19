package sacloud

import (
	"fmt"
	"github.com/sacloud/sackerel/job/core"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"time"
)

// CollectServerMetricsAllJob 過去分含めたサーバーメトリクス(CPU/Disk/Interface)を取得するジョブ
func CollectServerMetricsAllJob(payload interface{}) core.JobAPI {
	return core.NewJob("CollectMetricsServerAll", collectServerMetricsAll, payload)
}

// CollectServerMetricsLatestJob 直近のサーバーメトリクス(CPU/Disk/Interface)を取得するジョブ
func CollectServerMetricsLatestJob(payload interface{}) core.JobAPI {
	return core.NewJob("CollectMetricsServerLatest", collectServerMetricsLatest, payload)
}

func collectServerMetricsAll(queue *core.Queue, option *core.Option, job core.JobAPI) {
	collectServerMetricsInner(queue, option, job, nil)
}

func collectServerMetricsLatest(queue *core.Queue, option *core.Option, job core.JobAPI) {

	end := time.Now()
	start := end.Add(-option.MetricsHistoryPeriod)
	req := sacloud.NewResourceMonitorRequest(&start, &end)
	collectServerMetricsInner(queue, option, job, req)
}

func collectServerMetricsInner(queue *core.Queue, option *core.Option, job core.JobAPI, req *sacloud.ResourceMonitorRequest) {
	var payload = job.GetPayload()
	if payload == nil {
		queue.PushWarn(fmt.Errorf("'%s' => payload is nil", job.GetName()))
		return
	}

	// check payload types
	var server *sacloud.Server
	var sourcePayload *core.SourcePayload

	if s, ok := payload.(core.SourcePayloadHolder); ok {
		sourcePayload = s.GetSourcePayload()
		source := sourcePayload.SacloudSource
		if server, ok = source.(*sacloud.Server); !ok {
			queue.PushWarn(fmt.Errorf("'%s' => payload.Source is invalid type. need [*sacloud.Server]", job.GetName()))
			return
		}
	} else {
		queue.PushWarn(fmt.Errorf("'%s' => payload is invalid type. need [core.SourcePayloadHolder]", job.GetName()))
		return
	}

	// Create the collect-metrics payload
	metricsPayload := core.NewCollectMetricsPayload(sourcePayload)
	client := getClient(option, sourcePayload.SacloudZone)

	// cpu
	cpuMetrics, err := client.Server.Monitor(server.ID, req)
	if err != nil {
		queue.PushError(err)
		return
	}
	metricsPayload.Metrics.CPU = append(metricsPayload.Metrics.CPU, cpuMetrics)

	// disks
	for _, disk := range server.Disks {
		diskMetrics, err := client.Disk.Monitor(disk.ID, req)
		if err != nil {
			queue.PushError(err)
			return
		}
		metricsPayload.Metrics.Disk = append(metricsPayload.Metrics.Disk, diskMetrics)
	}

	// interfaces
	for _, nic := range server.Interfaces {
		nicMetrics, err := client.Interface.Monitor(nic.ID, req)
		if err != nil {
			queue.PushError(err)
			return
		}
		metricsPayload.Metrics.Interface = append(metricsPayload.Metrics.Interface, nicMetrics)
	}

	queue.PushRequest("collected-server-metrics", metricsPayload)
}
