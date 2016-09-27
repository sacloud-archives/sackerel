// sackerel - A integration tool of SakuraCloud and Mackerel.
package main

import (
	"fmt"
	"github.com/sacloud/sackerel/job"
	"github.com/sacloud/sackerel/job/core"
	"github.com/sacloud/sackerel/version"
	"gopkg.in/urfave/cli.v2"
	"io"
	"os"
	"strings"
)

var (
	appName              = "sackerel"
	appUsage             = "A integration tool of Mackerel and SakuraCloud"
	appCopyright         = "Copyright (C) 2016 Kazumichi Yamamoto."
	applHelpTextTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.HelpName}} [options]

REQUIRED PARAMETERS:
   {{range .RequiredFlags}}{{.}}
   {{end}}
OPTIONS:
   {{range .NormalFlags}}{{.}}
   {{end}}
************* ADVANCED OPTIONS **************

   FOR PERFORMANCE TUNING:
      {{range .PerformanceFlags}}{{.}}
      {{end}}
   FOR DEBUG:
      {{range .ForDeveloperFlags}}{{.}}
      {{end}}
*************************************************
VERSION:
   {{.Version}}

{{.Copyright}}
`

	requiredFlagNames = []string{
		"token",
		"secret",
		"apikey",
	}

	forDeveloperFlagNames = []string{
		"sakuracloud-trace-mode",
		"mackerel-trace-mode",
		"trace-log",
		"info-log",
		"warn-log",
		"error-log",
	}

	performanceFlagNames = []string{
		"reconcile-job-interval",
		"api-call-interval",
		"job-queue-size",
		"throttled-api-request-size",
		"sakura-api-queue-size",
		"mackerel-api-queue-size",
	}
)

func main() {

	// !!HACK!! disable HTTP/2 for MackerelAPI Server(NGINX)
	// 現状では、HTTP/2でPOSTするとエラーになるためここで回避しておく
	os.Setenv("GODEBUG", "http2client=0")

	cli.AppHelpTemplate = applHelpTextTemplate
	app := &cli.App{}
	option := core.NewOption()

	app.Name = appName
	app.Usage = appUsage
	app.HelpName = appName
	app.Copyright = appCopyright

	app.Flags = cliFlags(option)
	app.Action = cliCommand(option)
	app.Version = version.FullVersion()

	originalHelpPrinter := cli.HelpPrinter
	cli.HelpPrinter = func(w io.Writer, templ string, d interface{}) {
		app := d.(*cli.App)
		data := newHelpData(app)
		originalHelpPrinter(w, templ, data)
	}

	app.Run(os.Args)
}

type helpData struct {
	*cli.App
	RequiredFlags     []cli.Flag
	NormalFlags       []cli.Flag
	PerformanceFlags  []cli.Flag
	ForDeveloperFlags []cli.Flag
}

func newHelpData(app *cli.App) interface{} {
	data := &helpData{App: app}

	for _, f := range app.VisibleFlags() {
		if isExistsFlag(requiredFlagNames, f) {
			data.RequiredFlags = append(data.RequiredFlags, f)
		} else if isExistsFlag(forDeveloperFlagNames, f) {
			data.ForDeveloperFlags = append(data.ForDeveloperFlags, f)
		} else if isExistsFlag(performanceFlagNames, f) {
			data.PerformanceFlags = append(data.PerformanceFlags, f)
		} else {
			data.NormalFlags = append(data.NormalFlags, f)
		}
	}

	return data
}

func cliFlags(option *core.Option) []cli.Flag {

	return []cli.Flag{
		&cli.StringFlag{
			Name:        "token",
			Aliases:     []string{"sakuracloud-access-token"},
			Usage:       "API Token of SakuraCloud",
			EnvVars:     []string{"SAKURACLOUD_ACCESS_TOKEN"},
			DefaultText: "none",
			Destination: &option.SakuraCloudOption.AccessToken,
		},
		&cli.StringFlag{
			Name:        "secret",
			Aliases:     []string{"sakuracloud-access-token-secret"},
			Usage:       "API Secret of SakuraCloud",
			EnvVars:     []string{"SAKURACLOUD_ACCESS_TOKEN_SECRET"},
			DefaultText: "none",
			Destination: &option.SakuraCloudOption.AccessTokenSecret,
		},
		&cli.StringSliceFlag{
			Name:    "zones",
			Aliases: []string{"sakuracloud-zones"},
			Usage:   "Target zone list of SakuraCloud",
			EnvVars: []string{"SAKURACLOUD_ZONES"},
			Value:   cli.NewStringSlice("is1b", "tk1a"),
		},
		&cli.BoolFlag{
			Name:        "sakuracloud-trace-mode",
			Usage:       "Flag of SakuraCloud debug-mode",
			EnvVars:     []string{"SAKURACLOUD_TRACE_MODE"},
			Destination: &option.SakuraCloudOption.TraceMode,
			Value:       false,
		},
		&cli.StringFlag{
			Name:        "apikey",
			Aliases:     []string{"mackerel-apikey"},
			Usage:       "API Key of Mackerel",
			EnvVars:     []string{"MACKEREL_APIKEY"},
			DefaultText: "none",
			Destination: &option.MackerelOption.APIKey,
		},
		&cli.BoolFlag{
			Name:        "mackerel-trace-mode",
			Usage:       "Flag of Mackerel debug-mode",
			EnvVars:     []string{"MACKEREL_TRACE_MODE"},
			Destination: &option.MackerelOption.TraceMode,
			Value:       false,
		},
		&cli.DurationFlag{
			Name:        "interval",
			Aliases:     []string{"timer-job-interval"},
			Usage:       "Interval of each timer jobs",
			EnvVars:     []string{"SACKEREL_TIMER_JOB_INTERVAL"},
			Destination: &option.TimerJobInterval,
			Value:       core.DefaultTimerJobInterval,
		},
		&cli.DurationFlag{
			Name:        "reconcile-job-interval",
			Usage:       "Interval of each reconcile jobs",
			EnvVars:     []string{"SACKEREL_RECONCILE_JOB_INTERVAL"},
			Destination: &option.ReconcileJobInterval,
			Value:       core.DefaultReconcileJobInterval,
		},
		&cli.DurationFlag{
			Name:        "period",
			Aliases:     []string{"metrics-history-period"},
			Usage:       "Period of collecting metrics history",
			EnvVars:     []string{"SACKEREL_METRICS_HISTORY_PERIOD"},
			Destination: &option.MetricsHistoryPeriod,
			Value:       core.DefaultMetricsHistoryPeriod,
		},
		&cli.DurationFlag{
			Name:        "api-call-interval",
			Usage:       "Time duration of API call interval",
			EnvVars:     []string{"SACKEREL_API_CALL_INTERVAL"},
			Destination: &option.APICallInterval,
			Value:       core.DefaultAPICallInterval,
		},
		&cli.IntFlag{
			Name:        "port",
			Aliases:     []string{"healthcheck-port"},
			Usage:       "Number of web server port for healthcheck",
			EnvVars:     []string{"SACKEREL_HEALTHCHECK_PORT"},
			Destination: &option.HealthCheckWebServerPort,
			Value:       core.DefaultHealthCheckWebServerPort,
		},
		&cli.IntFlag{
			Name:        "job-queue-size",
			Usage:       "Size of internal job queue",
			EnvVars:     []string{"SACKEREL_JOB_QUEUE_SIZE"},
			Destination: &option.JobQueueBufSize,
			Value:       core.DefaultJobQueueBufSize,
		},
		&cli.IntFlag{
			Name:        "throttled-api-request-size",
			Usage:       "Size of throttledAPI requst queue",
			EnvVars:     []string{"SACKEREL_THROTTLED_API_REQUST_QUEUE_SIZE"},
			Destination: &option.ThrottledAPIReqQueueBufSize,
			Value:       core.DefaultThrottledAPIReqQueueBufSize,
		},
		&cli.IntFlag{
			Name:        "sakura-api-queue-size",
			Usage:       "Size of SauraAPI request queue",
			EnvVars:     []string{"SACKEREL_SAKURA_API_REQEST_QUEUE_SIZE"},
			Destination: &option.SakuraCloudOption.ReuqestQueueBufSize,
			Value:       core.DefaultSakuraAPIReqQueueBufSize,
		},
		&cli.IntFlag{
			Name:        "mackerel-api-queue-size",
			Usage:       "Size of MackerelAPI request queue",
			EnvVars:     []string{"SACKEREL_MACKEREL_API_REQEST_QUEUE_SIZE"},
			Destination: &option.MackerelOption.RequestQueueBufSize,
			Value:       core.DefaultMackerelAPIReqQueueBufSize,
		},
		&cli.BoolFlag{
			Name:        "disable-healthcheck",
			Usage:       "Flag of disable health check web server",
			EnvVars:     []string{"SACKEREL_DISABLE_HEALTHCHECK"},
			Destination: &option.DisableHealthCheck,
			Value:       false,
		},
		&cli.BoolFlag{
			Name:        "skip-init",
			Usage:       "Flag of skip init job",
			EnvVars:     []string{"SACKEREL_SKIP_INIT"},
			Destination: &option.SkipInit,
			Value:       false,
		},
		&cli.BoolFlag{
			Name:        "trace-log",
			Usage:       "Flag of enable TRACE log",
			EnvVars:     []string{"SACKEREL_TRACE_LOG"},
			Destination: &option.TraceLog,
			Value:       false,
		},
		&cli.BoolFlag{
			Name:        "info-log",
			Usage:       "Flag of enable INFO log",
			EnvVars:     []string{"SACKEREL_INFO_LOG"},
			Value:       true,
			Destination: &option.InfoLog,
		},
		&cli.BoolFlag{
			Name:        "warn-log",
			Usage:       "Flag of enable WARN log",
			EnvVars:     []string{"SACKEREL_WARN_LOG"},
			Value:       true,
			Destination: &option.WarnLog,
		},
		&cli.BoolFlag{
			Name:        "error-log",
			Usage:       "Flag of enable ERROR log",
			EnvVars:     []string{"SACKEREL_ERROR_LOG"},
			Value:       true,
			Destination: &option.ErrorLog,
		},
	}

}

func cliCommand(option *core.Option) func(c *cli.Context) error {
	return func(c *cli.Context) error {

		option.SakuraCloudOption.Zone = c.StringSlice("sakuracloud-zones")
		errors := option.Validate()
		if len(errors) != 0 {
			return flattenErrors(errors)
		}

		return job.Run(option)

	}
}

func flattenErrors(errors []error) error {
	var list = make([]string, 0)
	for _, str := range errors {
		list = append(list, str.Error())
	}
	return fmt.Errorf(strings.Join(list, "\n"))
}

func isExistsFlag(source []string, target cli.Flag) bool {
	for _, s := range source {
		if s == target.Names()[0] {
			return true
		}
	}
	return false
}
