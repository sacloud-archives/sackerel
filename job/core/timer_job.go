package core

import "time"

// TimerJobAPI タイマー起動されるジョブを表すインターフェイス
type TimerJobAPI interface {
	JobAPI
	GetTickerDuration() time.Duration
}

// TimerJob タイマー起動されるジョブ
type TimerJob struct {
	*Job
	TickerDuration time.Duration
}

// GetTickerDuration タイマー起動間隔の取得
func (j *TimerJob) GetTickerDuration() time.Duration {
	return j.TickerDuration
}

// NewTimerJob タイマー起動されるジョブの作成
func NewTimerJob(name string, workerFunc WorkerFunc, tickerDuration time.Duration) *TimerJob {
	return &TimerJob{
		Job:            NewJob(name, workerFunc, nil),
		TickerDuration: tickerDuration,
	}
}
