package mackerel

import (
	"github.com/sacloud/sackerel/job/core"
)

// FindServerJob Mackerel上のホストからサーバーを検索するジョブ
func FindServerJob(payload interface{}) core.JobAPI {
	return core.NewJob("MackerelFindServer", findMackerelServer, payload)
}

// FindRouterJob Mackerel上のホストからルーターを検索するジョブ
func FindRouterJob(payload interface{}) core.JobAPI {
	return core.NewJob("MackerelFindRouter", findMackerelRouter, payload)
}

// FindLoadBalancerJob Mackerel上のホストからロードバランサーを検索するジョブ
func FindLoadBalancerJob(payload interface{}) core.JobAPI {
	return core.NewJob("MackerelFindLoadBalancer", findMackerelLoadBalancer, payload)
}

// FindVPCRouterJob Mackerel上のホストからVPCルーターを検索するジョブ
func FindVPCRouterJob(payload interface{}) core.JobAPI {
	return core.NewJob("MackerelFindVPCRouter", findMackerelVPCRouter, payload)
}

// FindDatabaseJob Mackerel上のホストからデータベースを検索するジョブ
func FindDatabaseJob(payload interface{}) core.JobAPI {
	return core.NewJob("MackerelFindDatabase", findMackerelDatabase, payload)
}

func findMackerelServer(queue *core.Queue, option *core.Option, job core.JobAPI) {
	findMackerelTarget(queue, option, job, "server")
}
func findMackerelRouter(queue *core.Queue, option *core.Option, job core.JobAPI) {
	findMackerelTarget(queue, option, job, "router")
}
func findMackerelLoadBalancer(queue *core.Queue, option *core.Option, job core.JobAPI) {
	findMackerelTarget(queue, option, job, "loadbalancer")
}
func findMackerelVPCRouter(queue *core.Queue, option *core.Option, job core.JobAPI) {
	findMackerelTarget(queue, option, job, "vpcrouter")
}
func findMackerelDatabase(queue *core.Queue, option *core.Option, job core.JobAPI) {
	findMackerelTarget(queue, option, job, "database")
}
