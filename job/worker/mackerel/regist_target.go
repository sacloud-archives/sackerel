package mackerel

import (
	"github.com/sacloud/sackerel/job/core"
)

// RegistServerJob Mackerelホストとしてサーバーを登録するジョブ
func RegistServerJob(payload interface{}) core.JobAPI {
	return core.NewJob("MackerelRegistServer", registServer, payload)
}

// RegistRouterJob Mackerelホストとしてルーターを登録するジョブ
func RegistRouterJob(payload interface{}) core.JobAPI {
	return core.NewJob("MackerelRegistRouter", registRouter, payload)
}

// RegistLoadBalancerJob Mackerelホストとしてロードバランサーを登録するジョブ
func RegistLoadBalancerJob(payload interface{}) core.JobAPI {
	return core.NewJob("MackerelRegistLoadBalancer", registLoadBalancer, payload)
}

// RegistVPCRouterJob MackerelホストとしてVPCルーターを登録するジョブ
func RegistVPCRouterJob(payload interface{}) core.JobAPI {
	return core.NewJob("MackerelRegistVPCRouter", registVPCRouter, payload)
}

// RegistDatabaseJob Mackerelホストとしてデータベースを登録するジョブ
func RegistDatabaseJob(payload interface{}) core.JobAPI {
	return core.NewJob("MackerelRegistDatabase", registDatabase, payload)
}

func registServer(queue *core.Queue, option *core.Option, job core.JobAPI) {
	registMackerelTarget(queue, option, job, "server")
}
func registRouter(queue *core.Queue, option *core.Option, job core.JobAPI) {
	registMackerelTarget(queue, option, job, "router")
}
func registLoadBalancer(queue *core.Queue, option *core.Option, job core.JobAPI) {
	registMackerelTarget(queue, option, job, "loadbalancer")
}
func registVPCRouter(queue *core.Queue, option *core.Option, job core.JobAPI) {
	registMackerelTarget(queue, option, job, "vpcrouter")
}
func registDatabase(queue *core.Queue, option *core.Option, job core.JobAPI) {
	registMackerelTarget(queue, option, job, "database")
}
