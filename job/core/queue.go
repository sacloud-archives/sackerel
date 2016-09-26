package core

// Queue ジョブキュー
type Queue struct {
	Request          chan JobRequestAPI
	Internal         chan JobAPI
	SakuraRequest    chan JobAPI
	MackerelRequest  chan JobAPI
	ThrottledRequest chan JobAPI
	Logs             *LogQueue
	Quit             chan error
}

// LogQueue ログ出力用キュー
type LogQueue struct {
	Info  chan string
	Trace chan string
	Warn  chan error
	Error chan error
}

var defaultLogBufferSize = 10

// NewQueue ジョブキューの新規作成
func NewQueue(workBufSize int, throttledReqBufSize int, sakuraReqBufSize int, mkrReqBufSize int) *Queue {
	return &Queue{
		Request:          make(chan JobRequestAPI, workBufSize),
		Internal:         make(chan JobAPI, workBufSize),
		ThrottledRequest: make(chan JobAPI, throttledReqBufSize),
		SakuraRequest:    make(chan JobAPI, sakuraReqBufSize),
		MackerelRequest:  make(chan JobAPI, mkrReqBufSize),
		Logs: &LogQueue{
			Info:  make(chan string, defaultLogBufferSize),
			Trace: make(chan string, defaultLogBufferSize),
			Warn:  make(chan error, defaultLogBufferSize),
			Error: make(chan error, defaultLogBufferSize),
		},
		Quit: make(chan error),
	}
}

// PushRequest push new request to job-routing queue
func (q *Queue) PushRequest(requestName string, payload interface{}) {
	q.Request <- &jobRequest{
		name:    requestName,
		payload: payload,
	}
}

//---------------------------------------------------------
// Push jobs
//---------------------------------------------------------

// PushInternalWork push job to internal-job queue
func (q *Queue) PushInternalWork(work JobAPI) {
	q.Internal <- work
}

// PushSakuraAPIWork push job to sakuraAPI-job queue
func (q *Queue) PushSakuraAPIWork(work JobAPI) {
	q.SakuraRequest <- work
}

// PushMackerelAPIWork push job to mackerelAPI-job queue
func (q *Queue) PushMackerelAPIWork(work JobAPI) {
	q.MackerelRequest <- work
}

// PushThrottledAPIWork push job to throttledAPI-job queue
func (q *Queue) PushThrottledAPIWork(work JobAPI) {
	q.ThrottledRequest <- work
}

//---------------------------------------------------------
// Stop
//---------------------------------------------------------

// Stop push stop request to queue
func (q *Queue) Stop() {
	q.Quit <- nil
}

// StopByError push stop request wth error to queue
func (q *Queue) StopByError(err error) {
	q.Quit <- err
}

//---------------------------------------------------------
// Logging functions
//---------------------------------------------------------

// PushTrace push message to trace-log queue
func (q *Queue) PushTrace(msg string) {
	q.Logs.Trace <- msg
}

// PushInfo push message to info-log queue
func (q *Queue) PushInfo(msg string) {
	q.Logs.Info <- msg
}

// PushWarn push message to warn-log queue
func (q *Queue) PushWarn(err error) {
	q.Logs.Warn <- err
}

// PushError push error to error queue
func (q *Queue) PushError(err error) {
	q.Logs.Error <- err
}
