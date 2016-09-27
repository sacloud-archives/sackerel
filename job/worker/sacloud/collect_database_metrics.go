package sacloud

import (
	"fmt"
	"github.com/sacloud/sackerel/job/core"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"time"
)

// CollectDatabaseMetricsAllJob 過去分含めたデータベースメトリクス(cpu/memory/interface/disk)を取得するジョブ
func CollectDatabaseMetricsAllJob(payload interface{}) core.JobAPI {
	return core.NewJob("CollectMetricsDatabaseAll", collectDatabaseMetricsAll, payload)
}

// CollectDatabaseMetricsLatestJob 直近のデータベースメトリクス(Interface)を取得するジョブ
func CollectDatabaseMetricsLatestJob(payload interface{}) core.JobAPI {
	return core.NewJob("CollectMetricsDatabaseLatest", collectDatabaseMetricsLatest, payload)
}

func collectDatabaseMetricsAll(queue *core.Queue, option *core.Option, job core.JobAPI) {
	collectDatabaseMetricsInner(queue, option, job, nil)
}

func collectDatabaseMetricsLatest(queue *core.Queue, option *core.Option, job core.JobAPI) {

	end := time.Now()
	start := end.Add(-option.MetricsHistoryPeriod)
	req := sacloud.NewResourceMonitorRequest(&start, &end)
	collectDatabaseMetricsInner(queue, option, job, req)
}

func collectDatabaseMetricsInner(queue *core.Queue, option *core.Option, job core.JobAPI, req *sacloud.ResourceMonitorRequest) {
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
		if _, ok := source.(*sacloud.Database); !ok {
			queue.PushWarn(fmt.Errorf("'%s' => payload.Source is invalid type. need [*sacloud.Database]", job.GetName()))
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

	//cpu
	cpuMetrics, err := client.Database.MonitorCPU(sourcePayload.SacloudResourceID, req)
	if err != nil {
		queue.PushError(err)
		return
	}
	metricsPayload.Metrics.CPU = append(metricsPayload.Metrics.CPU, cpuMetrics)

	// interfaces
	nicMetrics, err := client.Database.MonitorInterface(sourcePayload.SacloudResourceID, req)
	if err != nil {
		queue.PushError(err)
		return
	}
	metricsPayload.Metrics.Interface = append(metricsPayload.Metrics.Interface, nicMetrics)

	// Disk:System(Read/Write)
	systemDiskMetrics, err := client.Database.MonitorSystemDisk(sourcePayload.SacloudResourceID, req)
	if err != nil {
		queue.PushError(err)
		return
	}
	metricsPayload.Metrics.Disk = append(metricsPayload.Metrics.Interface, systemDiskMetrics)

	// Disk:Backup(Read/Write)
	backupDiskMetrics, err := client.Database.MonitorBackupDisk(sourcePayload.SacloudResourceID, req)
	if err != nil {
		queue.PushError(err)
		return
	}
	metricsPayload.Metrics.Disk = append(metricsPayload.Metrics.Interface, backupDiskMetrics)

	// Database(MemorySize , DiskSize[System/Backup])
	databaseMetrics, err := client.Database.MonitorDatabase(sourcePayload.SacloudResourceID, req)
	if err != nil {
		queue.PushError(err)
		return
	}
	metricsPayload.Metrics.Database = append(metricsPayload.Metrics.Interface, databaseMetrics)

	queue.PushRequest("collected-database-metrics", metricsPayload)
}
