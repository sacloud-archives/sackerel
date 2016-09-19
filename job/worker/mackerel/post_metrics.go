package mackerel

import (
	"fmt"
	"github.com/sacloud/sackerel/job/core"
)

// PostMetricsJob メトリクス投入用ジョブ
func PostMetricsJob(payload interface{}) core.JobAPI {
	return core.NewJob("MackerelPostMetrics", postMetrics, payload)
}

func postMetrics(queue *core.Queue, option *core.Option, job core.JobAPI) {

	var payload = job.GetPayload()
	if payload == nil {
		queue.PushWarn(fmt.Errorf("'%s' => payload is nil", job.GetName()))
		return
	}

	if metricsPayload, ok := payload.(*core.CollectMetricsPayload); ok {

		if len(metricsPayload.MackerelMetricsParam) > 0 {
			client := getClient(option)
			err := client.PostHostMetricValues(metricsPayload.MackerelMetricsParam)

			if err != nil {
				queue.PushError(err)
				return
			}
		}

		targetName := fmt.Sprintf("(%s)", metricsPayload.SourceType.Name())
		queue.PushInfo(
			fmt.Sprintf(
				"Metrics posted %-15s => SakuraID:[%d] / MackerelID:[%s] / MetricsCount:[%d]",
				targetName,
				metricsPayload.SacloudResourceID,
				metricsPayload.MackerelID,
				len(metricsPayload.MackerelMetricsParam),
			),
		)

		queue.PushRequest("posted-metrics", payload)
	} else {
		queue.PushWarn(fmt.Errorf("'%s' => payload is invalid type. need [core.CollectMetricsPayload]", job.GetName()))
		return
	}

}
