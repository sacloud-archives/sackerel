package job

import (
	"fmt"
	"github.com/sacloud/sackerel/job/core"
)

// Router ジョブキューに対するルーティング定義/処理を提供する
type Router struct {
	queue  *core.Queue
	option *core.Option
	routes map[string]core.JobRouterFunc
}

// NewRouter Routerの新規作成
func NewRouter(queue *core.Queue, option *core.Option) *Router {
	r := &Router{
		queue:  queue,
		option: option,
	}
	r.buildRouteDefines()
	return r
}

// Routing ジョブキューへの処理リクエストを定義に従いルーティングし、適切なワーカーを呼び出す
func (r *Router) Routing(req core.JobRequestAPI) {

	payload := req.GetPayload()
	r.queue.PushTrace(fmt.Sprintf("request => '%s' payload => (%#v)", req.GetName(), payload))

	if route, ok := r.routes[req.GetName()]; ok {
		route(r.queue, r.option, req)
	} else {
		r.queue.PushWarn(fmt.Errorf("Route('%s') is not found.", req.GetName()))
	}
}
