package mackerel

import (
	mkr "github.com/mackerelio/mackerel-client-go"
	"github.com/sacloud/sackerel/job/core"
)

var (
	cpuCustomGraph = &mkr.GraphDefsParam{
		Name:        "custom.sacloud.cpu.#",
		DisplayName: "CPU",
		Unit:        "percentage",
		Metrics: []*mkr.GraphDefsMetric{
			{
				Name:        "custom.sacloud.cpu.#.cpu-time",
				DisplayName: "CPU Time",
				IsStacked:   false,
			},
		},
	}

	diskCustomGraph = &mkr.GraphDefsParam{
		Name:        "custom.sacloud.disk.#",
		DisplayName: "Disk",
		Unit:        "bytes/sec",
		Metrics: []*mkr.GraphDefsMetric{
			{
				Name:        "custom.sacloud.disk.#.read",
				DisplayName: "Read bytes",
				IsStacked:   false,
			},
			{
				Name:        "custom.sacloud.disk.#.write",
				DisplayName: "Write bytes",
				IsStacked:   false,
			},
		},
	}

	interfaceCustomGraph = &mkr.GraphDefsParam{
		Name:        "custom.sacloud.interface.#",
		DisplayName: "Interface",
		Unit:        "bytes/sec",
		Metrics: []*mkr.GraphDefsMetric{
			{
				Name:        "custom.sacloud.interface.#.send",
				DisplayName: "Send bytes",
				IsStacked:   false,
			},
			{
				Name:        "custom.sacloud.interface.#.receive",
				DisplayName: "Receive bytes",
				IsStacked:   false,
			},
		},
	}

	memorySizeCustomGraph = &mkr.GraphDefsParam{
		Name:        "custom.sacloud.memorysize.#",
		DisplayName: "Memory Size",
		Unit:        "bytes",
		Metrics: []*mkr.GraphDefsMetric{
			{
				Name:        "custom.sacloud.memorysize.#.total",
				DisplayName: "Total memory bytes",
				IsStacked:   false,
			},
			{
				Name:        "custom.sacloud.memorysize.#.used",
				DisplayName: "Used memory bytes",
				IsStacked:   false,
			},
		},
	}
	diskSizeBackupCustomGraph = &mkr.GraphDefsParam{
		Name:        "custom.sacloud.disksize.backup.#",
		DisplayName: "Disk Size(Backup)",
		Unit:        "bytes",
		Metrics: []*mkr.GraphDefsMetric{
			{
				Name:        "custom.sacloud.disksize.backup.#.total",
				DisplayName: "Total disk bytes",
				IsStacked:   false,
			},
			{
				Name:        "custom.sacloud.disksize.backup.#.used",
				DisplayName: "Used disk bytes",
				IsStacked:   false,
			},
		},
	}
	diskSizeSystemCustomGraph = &mkr.GraphDefsParam{
		Name:        "custom.sacloud.disksize.system.#",
		DisplayName: "Disk Size(System)",
		Unit:        "bytes",
		Metrics: []*mkr.GraphDefsMetric{
			{
				Name:        "custom.sacloud.disksize.system.#.total",
				DisplayName: "Total disk bytes",
				IsStacked:   false,
			},
			{
				Name:        "custom.sacloud.disksize.system.#.used",
				DisplayName: "Used disk bytes",
				IsStacked:   false,
			},
		},
	}
)

// RegistGraphDefsJob Mackerelのカスタムグラフ登録用ジョブ
func RegistGraphDefsJob(payload interface{}) core.JobAPI {
	return core.NewJob("MackerelRegistGraphDefs", registGraphDefs, payload)
}

func registGraphDefs(queue *core.Queue, option *core.Option, job core.JobAPI) {

	client := getClient(option)

	err := client.CreateGraphDefs([]*mkr.GraphDefsParam{
		cpuCustomGraph,
		diskCustomGraph,
		interfaceCustomGraph,
		memorySizeCustomGraph,
		diskSizeBackupCustomGraph,
		diskSizeSystemCustomGraph,
	})

	if err != nil {
		queue.PushError(err)
	}

	queue.PushInfo("Initialized custom graph defines")
}
