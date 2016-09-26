package core

import "fmt"

// WorkerFunc ジョブキューでのワーカーが保有する処理用関数シグニチャ
type WorkerFunc func(*Queue, *Option, JobAPI)

// JobAPI ジョブキューでのジョブ(ワーカー呼び出し)を表すインターフェース
type JobAPI interface {
	GetName() string
	GetPayload() interface{}
	Start(*Queue, *Option)
}

// Job ジョブキューでのジョブ定義
type Job struct {
	Name       string
	WorkerFunc WorkerFunc
	Payload    interface{}
}

// Start ジョブ内のワーカー呼び出し
func (w *Job) Start(queue *Queue, option *Option) {
	if w.WorkerFunc == nil {
		queue.PushError(fmt.Errorf("'%s': WorkerFunc is requied", w.GetName()))
		return
	}

	w.WorkerFunc(queue, option, w)
}

// GetName ジョブ名取得
func (w *Job) GetName() string {
	return w.Name
}

// GetPayload ペイロード取得
func (w *Job) GetPayload() interface{} {
	return w.Payload
}

// NewJob Jobの新規作成
func NewJob(name string, workerFunc WorkerFunc, payload interface{}) *Job {
	return &Job{
		Name:       name,
		WorkerFunc: workerFunc,
		Payload:    payload,
	}
}
