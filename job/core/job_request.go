package core

// JobRouterFunc ジョブキューでのルーティング処理用関数シグニチャ
type JobRouterFunc func(*Queue, *Option, JobRequestAPI)

// JobRequestAPI ジョブキューへのルーティング要求時パラメータ
type JobRequestAPI interface {
	GetName() string
	GetPayload() interface{}
}

type jobRequest struct {
	name    string
	payload interface{}
}

// GetName ジョブ名称の取得
func (w *jobRequest) GetName() string {
	return w.name
}

// GetPayload ペイロードの取得
func (w *jobRequest) GetPayload() interface{} {
	return w.payload
}
