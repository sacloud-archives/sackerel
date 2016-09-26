package exchange

import (
	"github.com/sacloud/sackerel/job/core"
	"github.com/stretchr/testify/assert"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"reflect"
	"testing"
)

var server = &sacloud.Server{
	Resource: &sacloud.Resource{ID: 123456789012},
	Name:     "Test",
	TagsType: &sacloud.TagsType{
		Tags: []string{"tag1", "tag2"},
	},
	Interfaces: []sacloud.Interface{
		{
			IPAddress:  "192.168.0.1",
			MACAddress: "00:00::00:00:00:00:01",
		},
		{
			IPAddress:  "192.168.0.2",
			MACAddress: "00:00::00:00:00:00:02",
		},
	},
}

func TestServerJob(t *testing.T) {
	if !baseExchangeJobTest(t) {
		return
	}

	serverPayload := core.NewCreateHostPayload(server, "is1b", 123456789012, reflect.TypeOf(server))
	job := ServerJob(serverPayload)
	errFunc := func(err error) {
		assert.Fail(t, "Warn or Error log printed")
	}
	go job.Start(queue, option)
	jobCheckFunc(t, queue, func(v core.JobRequestAPI) {
		payload := v.GetPayload().(*core.CreateHostPayload)
		assert.NotNil(t, payload)

		para := payload.MackerelHostParam
		assert.NotNil(t, para)

		assert.Equal(t, para.CustomIdentifier, payload.GenerateMackerelName())
		assert.Equal(t, para.DisplayName, "Test")
		assert.Equal(t, para.Name, payload.GenerateMackerelName())
		assert.Equal(t, para.RoleFullnames, []string{
			"SakuraCloud:Server",
			"SakuraCloud:Zone-is1b",
			"SakuraCloud:tag1",
			"SakuraCloud:tag2",
		})
		assert.Equal(t, para.Interfaces[0].Name, "eth0")
		assert.Equal(t, para.Interfaces[0].IPAddress, "192.168.0.1")
		assert.Equal(t, para.Interfaces[0].MacAddress, "00:00::00:00:00:00:01")

		assert.Equal(t, para.Interfaces[1].Name, "eth1")
		assert.Equal(t, para.Interfaces[1].IPAddress, "192.168.0.2")
		assert.Equal(t, para.Interfaces[1].MacAddress, "00:00::00:00:00:00:02")

	}, errFunc, errFunc)
}
