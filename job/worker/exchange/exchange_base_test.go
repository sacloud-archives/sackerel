package exchange

import (
	"github.com/sacloud/sackerel/job/core"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

var queue = core.NewQueue(5, 5, 5, 5)
var option = core.NewOption()
var jobCheckFunc = func(t *testing.T, q *core.Queue, f func(v core.JobRequestAPI), w func(err error), e func(err error)) {
	for {
		select {
		case v := <-q.Request:
			if f != nil {
				f(v)
			}
			return
		case v := <-q.Logs.Warn:
			if w != nil {
				w(v)
			}
			return
		case v := <-q.Logs.Error:
			if e != nil {
				e(v)
			}
			return
		case <-time.After(5 * time.Second):
			assert.Fail(t, "Timeout")
			return
		}
	}
}

func baseExchangeJobTest(t *testing.T) bool {
	var job core.JobAPI
	// PayloadがnilだとWarn
	job = DatabaseJob(nil)
	go job.Start(queue, option)

	jobCheckFunc(t, queue, nil, func(err error) {
		assert.Error(t, err)
	}, nil)

	// PayloadにCreateHostPayload以外を渡すとWarn
	job = DatabaseJob("test")
	go job.Start(queue, option)
	jobCheckFunc(t, queue, nil, func(err error) {
		assert.Error(t, err)
	}, nil)

	// ソースに必要な型が設定されていなければエラー
	payload := core.NewCreateHostPayload("test", "is1b", 0, reflect.TypeOf("test"))
	job = DatabaseJob(payload)
	go job.Start(queue, option)
	jobCheckFunc(t, queue, nil, func(err error) {
		assert.Error(t, err)
	}, nil)

	return true
}
