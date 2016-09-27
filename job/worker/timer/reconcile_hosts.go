package timer

import (
	"github.com/sacloud/sackerel/job/core"
	"time"
)

// ReconcileHostTimerJob タイマー起動でのホストリコンサイルジョブ
func ReconcileHostTimerJob(timerDuration time.Duration) *core.TimerJob {
	return core.NewTimerJob("ReconsileAll", reconsileAll, timerDuration)
}

func reconsileAll(queue *core.Queue, option *core.Option, job core.JobAPI) {
	queue.PushRequest("reconcile-all", nil)
}
