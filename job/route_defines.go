package job

import (
	"fmt"
	"github.com/sacloud/sackerel/job/core"
	"github.com/sacloud/sackerel/job/worker/exchange"
	"github.com/sacloud/sackerel/job/worker/mackerel"
	"github.com/sacloud/sackerel/job/worker/sacloud"
)

func (r *Router) buildRouteDefines() {
	r.routes = map[string]core.JobRouterFunc{
		//--------------------
		// init
		//--------------------
		"init":       _PathThrough("init-all"),
		"init-all":   _Parallel([]string{"init-graph"}),
		"init-graph": _MackerelAPIRequest(mackerel.RegistGraphDefsJob),

		//--------------------
		// detect resources
		//--------------------
		"detect-resource-all": _Parallel(
			[]string{
				"detect-server",
				"detect-router",
				"detect-loadbalancer",
				"detect-vpcrouter",
				"detect-database",
			},
		),
		"detect-server":       _SakuraAPIRequest(sacloud.DetectServerJob),
		"detect-router":       _SakuraAPIRequest(sacloud.DetectRouterJob),
		"detect-loadbalancer": _SakuraAPIRequest(sacloud.DetectLoadBalancerJob),
		"detect-vpcrouter":    _SakuraAPIRequest(sacloud.DetectVPCRouterJob),
		"detect-database":     _SakuraAPIRequest(sacloud.DetectDatabaseJob),

		//--------------------
		// found-resource
		//--------------------
		"found-server":       _PathThrough("exchange-server"),
		"found-router":       _PathThrough("exchange-router"),
		"found-loadbalancer": _PathThrough("exchange-loadbalancer"),
		"found-vpcrouter":    _PathThrough("exchange-vpcrouter"),
		"found-database":     _PathThrough("exchange-database"),

		//--------------------
		// exchange
		//--------------------
		"exchange-server":       _ExchangeRequest(exchange.ServerJob),
		"exchange-router":       _ExchangeRequest(exchange.RouterJob),
		"exchange-loadbalancer": _ExchangeRequest(exchange.LoadBalancerJob),
		"exchange-vpcrouter":    _ExchangeRequest(exchange.VPCRouterJob),
		"exchange-database":     _ExchangeRequest(exchange.DatabaseJob),

		//--------------------
		// exchanged
		//--------------------
		"exchanged-server":       _PathThrough("mackerel-find-server"),
		"exchanged-router":       _PathThrough("mackerel-find-router"),
		"exchanged-loadbalancer": _PathThrough("mackerel-find-loadbalancer"),
		"exchanged-vpcrouter":    _PathThrough("mackerel-find-vpcrouter"),
		"exchanged-database":     _PathThrough("mackerel-find-database"),

		//--------------------
		// find from mackerel
		//--------------------
		"mackerel-find-server":       _MackerelAPIRequest(mackerel.FindServerJob),
		"mackerel-find-router":       _MackerelAPIRequest(mackerel.FindRouterJob),
		"mackerel-find-loadbalancer": _MackerelAPIRequest(mackerel.FindLoadBalancerJob),
		"mackerel-find-vpcrouter":    _MackerelAPIRequest(mackerel.FindVPCRouterJob),
		"mackerel-find-database":     _MackerelAPIRequest(mackerel.FindDatabaseJob),

		//--------------------
		// after find
		//--------------------
		"mackerel-found-server":       _Parallel([]string{"mackerel-update", "collect-server-metrics"}),
		"mackerel-found-router":       _Parallel([]string{"mackerel-update", "collect-router-metrics"}),
		"mackerel-found-loadbalancer": _Parallel([]string{"mackerel-update", "collect-loadbalancer-metrics"}),
		"mackerel-found-vpcrouter":    _Parallel([]string{"mackerel-update", "collect-vpcrouter-metrics"}),
		"mackerel-found-database":     _Parallel([]string{"mackerel-update", "collect-database-metrics"}),

		"mackerel-not-found-server":       _PathThrough("mackerel-regist-server"),
		"mackerel-not-found-router":       _PathThrough("mackerel-regist-router"),
		"mackerel-not-found-loadbalancer": _PathThrough("mackerel-regist-loadbalancer"),
		"mackerel-not-found-vpcrouter":    _PathThrough("mackerel-regist-vpcrouter"),
		"mackerel-not-found-database":     _PathThrough("mackerel-regist-database"),

		//--------------------
		// regist to mackerel
		//--------------------
		"mackerel-regist-server":       _ThrottledAPIRequest(mackerel.RegistServerJob),
		"mackerel-regist-router":       _ThrottledAPIRequest(mackerel.RegistRouterJob),
		"mackerel-regist-loadbalancer": _ThrottledAPIRequest(mackerel.RegistLoadBalancerJob),
		"mackerel-regist-vpcrouter":    _ThrottledAPIRequest(mackerel.RegistVPCRouterJob),
		"mackerel-regist-database":     _ThrottledAPIRequest(mackerel.RegistDatabaseJob),

		//--------------------
		// after regist
		//--------------------
		"mackerel-registed-server":       _Parallel([]string{"mackerel-update-status", "collect-server-metrics-all"}),
		"mackerel-registed-router":       _Parallel([]string{"mackerel-update-status", "collect-router-metrics-all"}),
		"mackerel-registed-loadbalancer": _Parallel([]string{"mackerel-update-status", "collect-loadbalancer-metrics-all"}),
		"mackerel-registed-vpcrouter":    _Parallel([]string{"mackerel-update-status", "collect-vpcrouter-metrics-all"}),
		"mackerel-registed-database":     _Parallel([]string{"mackerel-update-status", "collect-database-metrics-all"}),

		//--------------------
		// update mackerel
		//--------------------
		"mackerel-update": _Parallel([]string{"mackerel-update-host", "mackerel-update-status", "mackerel-update-role"}),

		//--------------------
		// update host
		//--------------------
		"mackerel-update-host":   _ThrottledAPIRequest(mackerel.UpdateHostJob),
		"mackerel-update-status": _ThrottledAPIRequest(mackerel.UpdateStatusJob),
		"mackerel-update-role":   _ThrottledAPIRequest(mackerel.UpdateRoleJob),

		//--------------------
		// after update host
		//--------------------
		"mackerel-updated-host":   _EndOfRoute,
		"mackerel-updated-status": _EndOfRoute,
		"mackerel-updated-role":   _EndOfRoute,

		//--------------------
		// collect
		//--------------------
		"collect-server-metrics":           _MackerelAPIRequest(sacloud.CollectServerMetricsLatestJob),
		"collect-server-metrics-all":       _MackerelAPIRequest(sacloud.CollectServerMetricsAllJob),
		"collect-router-metrics":           _MackerelAPIRequest(sacloud.CollectRouterMetricsLatestJob),
		"collect-router-metrics-all":       _MackerelAPIRequest(sacloud.CollectRouterMetricsAllJob),
		"collect-loadbalancer-metrics":     _MackerelAPIRequest(sacloud.CollectLoadBalancerMetricsLatestJob),
		"collect-loadbalancer-metrics-all": _MackerelAPIRequest(sacloud.CollectLoadBalancerMetricsAllJob),
		"collect-vpcrouter-metrics":        _MackerelAPIRequest(sacloud.CollectVPCRouterMetricsLatestJob),
		"collect-vpcrouter-metrics-all":    _MackerelAPIRequest(sacloud.CollectVPCRouterMetricsAllJob),
		"collect-database-metrics":         _MackerelAPIRequest(sacloud.CollectDatabaseMetricsLatestJob),
		"collect-database-metrics-all":     _MackerelAPIRequest(sacloud.CollectDatabaseMetricsAllJob),

		//--------------------
		// after collect
		//--------------------
		"collected-server-metrics":       _PathThrough("exchange-metrics"),
		"collected-router-metrics":       _PathThrough("exchange-metrics"),
		"collected-loadbalancer-metrics": _PathThrough("exchange-metrics"),
		"collected-vpcrouter-metrics":    _PathThrough("exchange-metrics"),
		"collected-database-metrics":     _PathThrough("exchange-metrics"),

		//--------------------
		// exchange - post
		//--------------------
		"exchange-metrics": _ExchangeRequest(exchange.MetricsJob),
		"post-metrics":     _MackerelAPIRequest(mackerel.PostMetricsJob),
		"posted-metrics":   _EndOfRoute,

		//--------------------
		// reconcile hosts
		//--------------------
		"reconcile-all":              _MackerelAPIRequest(mackerel.ReconcileByMACAddressJob),
		"reconcile-host":             _Parallel([]string{"add-agent-tag", "mackerel-retire-host"}),
		"add-agent-tag":              _SakuraAPIRequest(sacloud.AddAgentTagJob),
		"added-agent-tag":            _EndOfRoute,
		"mackerel-retire-host":       _ThrottledAPIRequest(mackerel.RetireJob),
		"mackerel-retired-host":      _PathThrough("mackerel-update-custom-id"),
		"mackerel-update-custom-id":  _ThrottledAPIRequest(mackerel.UpdadteCustomIDJob),
		"mackerel-updated-custom-id": _EndOfRoute,
	}
}

func _Parallel(routes []string) core.JobRouterFunc {
	return func(queue *core.Queue, option *core.Option, req core.JobRequestAPI) {
		for _, route := range routes {
			queue.PushRequest(route, req.GetPayload())
		}
	}
}

func _PathThrough(dest string) core.JobRouterFunc {
	return func(queue *core.Queue, option *core.Option, req core.JobRequestAPI) {
		queue.PushRequest(dest, req.GetPayload())
	}
}

func _ExchangeRequest(f func(interface{}) core.JobAPI) core.JobRouterFunc {
	return func(queue *core.Queue, option *core.Option, req core.JobRequestAPI) {
		queue.PushInternalWork(f(req.GetPayload()))
	}
}

func _SakuraAPIRequest(f func(interface{}) core.JobAPI) core.JobRouterFunc {
	return func(queue *core.Queue, option *core.Option, req core.JobRequestAPI) {
		queue.PushSakuraAPIWork(f(req.GetPayload()))
	}
}

func _MackerelAPIRequest(f func(interface{}) core.JobAPI) core.JobRouterFunc {
	return func(queue *core.Queue, option *core.Option, req core.JobRequestAPI) {
		queue.PushMackerelAPIWork(f(req.GetPayload()))
	}
}

func _ThrottledAPIRequest(f func(interface{}) core.JobAPI) core.JobRouterFunc {
	return func(queue *core.Queue, option *core.Option, req core.JobRequestAPI) {
		queue.PushThrottledAPIWork(f(req.GetPayload()))
	}
}

var _EndOfRoute core.JobRouterFunc = func(queue *core.Queue, option *core.Option, req core.JobRequestAPI) {
	queue.PushTrace(fmt.Sprintf("Route('%s') is finished.", req.GetName()))
}
