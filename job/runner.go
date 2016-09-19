package job

import (
	"fmt"
	"github.com/sacloud/sackerel/job/core"
)

// Run sackerelメイン処理
func Run(option *core.Option) error {

	// setup jobs environments
	jobQueue := core.NewQueue(option.QueueBufSizes())
	dispatcher := NewDispatcher(option, jobQueue)

	// start HealthCheck WebServer
	if !option.DisableHealthCheck {
		healthCheckServer := NewHealthCheckServer(option.HealthCheckWebServerPort)
		healthCheckServer.ListenAndServe()
		jobQueue.PushInfo(fmt.Sprintf("HealthCheck WebServer started on [:%d]", option.HealthCheckWebServerPort))
	}

	// start jobs
	err := dispatcher.Dispatch()

	if err != nil {
		return err
	}

	return nil
}
