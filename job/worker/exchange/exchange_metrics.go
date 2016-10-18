package exchange

import (
	"fmt"
	mkr "github.com/mackerelio/mackerel-client-go"
	"github.com/sacloud/libsacloud/sacloud"
	"github.com/sacloud/sackerel/job/core"
	"strings"
)

// MetricsJob メトリクス変換用ジョブ
func MetricsJob(payload interface{}) core.JobAPI {
	return core.NewJob("MackerelExchangeMetrics", exchangeMetrics, payload)
}

func exchangeMetrics(queue *core.Queue, option *core.Option, job core.JobAPI) {

	var payload = job.GetPayload()
	if payload == nil {
		queue.PushWarn(fmt.Errorf("'%s' => payload is nil", job.GetName()))
		return
	}

	if mPayload, ok := payload.(*core.CollectMetricsPayload); ok {

		exchangeFuncs := []func(*core.CollectMetricsPayload) error{
			exchangeCPUTimeMetrics,
			exchangeDiskReadMetrics,
			exchangeDiskWriteMetrics,
			exchangePacketReceiveMetrics,
			exchangePacketSendMetrics,
			exchangeInternetInMetrics,
			exchangeInternetOutMetrics,
			exchangeTotalMemorySizeMetrics,
			exchangeUsedMemorySizeMetrics,
			exchangeTotalDisk1SizeMetrics,
			exchangeUsedDisk1SizeMetrics,
			exchangeTotalDisk2SizeMetrics,
			exchangeUsedDisk2SizeMetrics,
		}

		for _, f := range exchangeFuncs {
			err := f(mPayload)
			if err != nil {
				queue.PushError(err)
				return
			}
		}

		queue.PushRequest("post-metrics", payload)

	} else {

		queue.PushWarn(fmt.Errorf("'%s' => payload is invalid type. need [core.CollectMetricsPayload]", job.GetName()))
		return
	}

}

func exchangeCPUTimeMetrics(mPayload *core.CollectMetricsPayload) error {

	metrics, err := extractMetricsParams(
		mPayload.MackerelID,
		"custom.sacloud.cpu.%d.cpu-time",
		mPayload.Metrics.CPU,
		extractCPUTimeValue)
	if err != nil {
		return err
	}
	mPayload.MackerelMetricsParam = append(mPayload.MackerelMetricsParam, metrics...)
	return nil
}

func exchangeDiskWriteMetrics(mPayload *core.CollectMetricsPayload) error {

	metrics, err := extractMetricsParams(
		mPayload.MackerelID,
		"custom.sacloud.disk.%d.write",
		mPayload.Metrics.Disk,
		extractDiskWriteValue)
	if err != nil {
		return err
	}
	mPayload.MackerelMetricsParam = append(mPayload.MackerelMetricsParam, metrics...)
	return nil
}

func exchangeDiskReadMetrics(mPayload *core.CollectMetricsPayload) error {

	metrics, err := extractMetricsParams(
		mPayload.MackerelID,
		"custom.sacloud.disk.%d.read",
		mPayload.Metrics.Disk,
		extractDiskReadValue)
	if err != nil {
		return err
	}
	mPayload.MackerelMetricsParam = append(mPayload.MackerelMetricsParam, metrics...)
	return nil
}

func exchangePacketSendMetrics(mPayload *core.CollectMetricsPayload) error {

	metrics, err := extractMetricsParams(
		mPayload.MackerelID,
		"custom.sacloud.interface.%d.send",
		mPayload.Metrics.Interface,
		extractPacketSendValue)
	if err != nil {
		return err
	}
	mPayload.MackerelMetricsParam = append(mPayload.MackerelMetricsParam, metrics...)
	return nil
}

func exchangePacketReceiveMetrics(mPayload *core.CollectMetricsPayload) error {

	metrics, err := extractMetricsParams(
		mPayload.MackerelID,
		"custom.sacloud.interface.%d.receive",
		mPayload.Metrics.Interface,
		extractPacketReceiveValue)
	if err != nil {
		return err
	}
	mPayload.MackerelMetricsParam = append(mPayload.MackerelMetricsParam, metrics...)
	return nil
}

func exchangeInternetInMetrics(mPayload *core.CollectMetricsPayload) error {

	metrics, err := extractMetricsParams(
		mPayload.MackerelID,
		"custom.sacloud.interface.%d.receive",
		mPayload.Metrics.Interface,
		extractInternetInValue)
	if err != nil {
		return err
	}
	mPayload.MackerelMetricsParam = append(mPayload.MackerelMetricsParam, metrics...)
	return nil
}

func exchangeInternetOutMetrics(mPayload *core.CollectMetricsPayload) error {

	metrics, err := extractMetricsParams(
		mPayload.MackerelID,
		"custom.sacloud.interface.%d.send",
		mPayload.Metrics.Interface,
		extractInternetOutValue)
	if err != nil {
		return err
	}
	mPayload.MackerelMetricsParam = append(mPayload.MackerelMetricsParam, metrics...)
	return nil
}

func exchangeTotalMemorySizeMetrics(mPayload *core.CollectMetricsPayload) error {

	metrics, err := extractMetricsParams(
		mPayload.MackerelID,
		"custom.sacloud.memorysize.%d.total",
		mPayload.Metrics.Database,
		extractTotalMemorySizeValue)
	if err != nil {
		return err
	}
	mPayload.MackerelMetricsParam = append(mPayload.MackerelMetricsParam, metrics...)
	return nil
}

func exchangeUsedMemorySizeMetrics(mPayload *core.CollectMetricsPayload) error {

	metrics, err := extractMetricsParams(
		mPayload.MackerelID,
		"custom.sacloud.memorysize.%d.used",
		mPayload.Metrics.Database,
		extractUsedMemorySizeValue)
	if err != nil {
		return err
	}
	mPayload.MackerelMetricsParam = append(mPayload.MackerelMetricsParam, metrics...)
	return nil
}

func exchangeTotalDisk1SizeMetrics(mPayload *core.CollectMetricsPayload) error {

	metrics, err := extractMetricsParams(
		mPayload.MackerelID,
		"custom.sacloud.disksize.system.%d.total",
		mPayload.Metrics.Database,
		extractTotalDisk1SizeValue)
	if err != nil {
		return err
	}
	mPayload.MackerelMetricsParam = append(mPayload.MackerelMetricsParam, metrics...)
	return nil
}
func exchangeUsedDisk1SizeMetrics(mPayload *core.CollectMetricsPayload) error {

	metrics, err := extractMetricsParams(
		mPayload.MackerelID,
		"custom.sacloud.disksize.system.%d.used",
		mPayload.Metrics.Database,
		extractUsedDisk1SizeValue)
	if err != nil {
		return err
	}
	mPayload.MackerelMetricsParam = append(mPayload.MackerelMetricsParam, metrics...)
	return nil
}

func exchangeTotalDisk2SizeMetrics(mPayload *core.CollectMetricsPayload) error {

	metrics, err := extractMetricsParams(
		mPayload.MackerelID,
		"custom.sacloud.disksize.backup.%d.total",
		mPayload.Metrics.Database,
		extractTotalDisk2SizeValue)
	if err != nil {
		return err
	}
	mPayload.MackerelMetricsParam = append(mPayload.MackerelMetricsParam, metrics...)
	return nil
}
func exchangeUsedDisk2SizeMetrics(mPayload *core.CollectMetricsPayload) error {

	metrics, err := extractMetricsParams(
		mPayload.MackerelID,
		"custom.sacloud.disksize.backup.%d.used",
		mPayload.Metrics.Database,
		extractUsedDisk2SizeValue)
	if err != nil {
		return err
	}
	mPayload.MackerelMetricsParam = append(mPayload.MackerelMetricsParam, metrics...)
	return nil
}

func extractCPUTimeValue(source *sacloud.MonitorValues) ([]sacloud.FlatMonitorValue, error) {
	values, err := source.FlattenCPUTimeValue()

	for i, v := range values {
		values[i].Value = v.Value * 100 //percentage
	}

	return values, err

}
func extractDiskWriteValue(source *sacloud.MonitorValues) ([]sacloud.FlatMonitorValue, error) {
	values, err := source.FlattenDiskWriteValue()
	return values, err
}
func extractDiskReadValue(source *sacloud.MonitorValues) ([]sacloud.FlatMonitorValue, error) {
	values, err := source.FlattenDiskReadValue()
	return values, err

}

func extractPacketSendValue(source *sacloud.MonitorValues) ([]sacloud.FlatMonitorValue, error) {
	values, err := source.FlattenPacketSendValue()
	return values, err
}

func extractPacketReceiveValue(source *sacloud.MonitorValues) ([]sacloud.FlatMonitorValue, error) {
	values, err := source.FlattenPacketReceiveValue()
	return values, err
}

func extractInternetInValue(source *sacloud.MonitorValues) ([]sacloud.FlatMonitorValue, error) {
	values, err := source.FlattenInternetInValue()
	for i, v := range values {
		if v.Value > 0 {
			values[i].Value = v.Value / 8 // bps to bytes/sec
		}
	}
	return values, err
}

func extractInternetOutValue(source *sacloud.MonitorValues) ([]sacloud.FlatMonitorValue, error) {
	values, err := source.FlattenInternetOutValue()
	for i, v := range values {
		if v.Value > 0 {
			values[i].Value = v.Value / 8 // bps to bytes/sec
		}
	}
	return values, err
}

func extractTotalMemorySizeValue(source *sacloud.MonitorValues) ([]sacloud.FlatMonitorValue, error) {
	values, err := source.FlattenTotalMemorySizeValue()
	for i, v := range values {
		if v.Value > 0 {
			values[i].Value = v.Value * 1024 // KB to byte
		}
	}
	return values, err
}

func extractUsedMemorySizeValue(source *sacloud.MonitorValues) ([]sacloud.FlatMonitorValue, error) {
	values, err := source.FlattenUsedMemorySizeValue()
	for i, v := range values {
		if v.Value > 0 {
			values[i].Value = v.Value * 1024 // KB to byte
		}
	}
	return values, err
}
func extractTotalDisk1SizeValue(source *sacloud.MonitorValues) ([]sacloud.FlatMonitorValue, error) {
	values, err := source.FlattenTotalDisk1SizeValue()
	for i, v := range values {
		if v.Value > 0 {
			values[i].Value = v.Value * 1024 // KB to byte
		}
	}
	return values, err
}
func extractUsedDisk1SizeValue(source *sacloud.MonitorValues) ([]sacloud.FlatMonitorValue, error) {
	values, err := source.FlattenUsedDisk1SizeValue()
	for i, v := range values {
		if v.Value > 0 {
			values[i].Value = v.Value * 1024 // KB to byte
		}
	}
	return values, err
}
func extractTotalDisk2SizeValue(source *sacloud.MonitorValues) ([]sacloud.FlatMonitorValue, error) {
	values, err := source.FlattenTotalDisk2SizeValue()
	for i, v := range values {
		if v.Value > 0 {
			values[i].Value = v.Value * 1024 // KB to byte
		}
	}
	return values, err
}
func extractUsedDisk2SizeValue(source *sacloud.MonitorValues) ([]sacloud.FlatMonitorValue, error) {
	values, err := source.FlattenUsedDisk2SizeValue()
	for i, v := range values {
		if v.Value > 0 {
			values[i].Value = v.Value * 1024 // KB to byte
		}
	}
	return values, err
}

func extractMetricsParams(mackerelID string, metricsKeyFormat string, from []*sacloud.MonitorValues, extractFunc func(*sacloud.MonitorValues) ([]sacloud.FlatMonitorValue, error)) ([]*mkr.HostMetricValue, error) {

	res := []*mkr.HostMetricValue{}
	for i, metrics := range from {
		if metrics == nil {
			continue
		}
		values, err := extractFunc(metrics)

		if err != nil {
			return res, err
		}

		// metricsKeyFormatにインデックス置き換え用プレースホルダが含まれる場合のみformat処理
		metricsLabel := metricsKeyFormat
		if strings.Contains(metricsLabel, "%d") {
			metricsLabel = fmt.Sprintf(metricsKeyFormat, i)
		}

		for _, v := range values {
			m := &mkr.HostMetricValue{
				HostID: mackerelID,
				MetricValue: &mkr.MetricValue{
					Name:  metricsLabel,
					Time:  v.Time.Unix(),
					Value: v.Value,
				},
			}

			res = append(res, m)
		}
	}
	return res, nil
}
