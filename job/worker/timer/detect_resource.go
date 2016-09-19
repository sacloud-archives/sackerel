package timer

import (
	"github.com/sacloud/sackerel/job/core"
	"time"
)

// DetectResourceTimerJob タイマー起動でのリソース情報収集ジョブ
func DetectResourceTimerJob(timerDuration time.Duration) *core.TimerJob {
	return core.NewTimerJob("DetectResource", startDetectResource, timerDuration)
}

func startDetectResource(queue *core.Queue, option *core.Option, job core.JobAPI) {
	queue.PushRequest("detect-resource-all", nil)
}
