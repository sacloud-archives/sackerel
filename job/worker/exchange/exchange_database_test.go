package exchange

import (
	"github.com/sacloud/sackerel/job/core"
	"github.com/stretchr/testify/assert"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"reflect"
	"testing"
)

var database = &sacloud.Database{
	Appliance: &sacloud.Appliance{
		Resource: &sacloud.Resource{ID: 123456789012},
		Name:     "Test",
		TagsType: &sacloud.TagsType{
			Tags: []string{"tag1", "tag2"},
		},
		Interfaces: []sacloud.Interface{
			{
				IPAddress: "8.8.8.8",
				Switch: &sacloud.Switch{
					Scope: sacloud.ESCopeShared,
				},
			},
		},
	},
	Remark: &sacloud.DatabaseRemark{
		ApplianceRemarkBase: &sacloud.ApplianceRemarkBase{
			Servers: []interface{}{
				map[string]interface{}{
					"IPAddress": "192.168.0.1",
				},
			},
		},
	},
}

func TestDatabaseJob(t *testing.T) {
	if !baseExchangeJobTest(t) {
		return
	}

	databasePayload := core.NewCreateHostPayload(database, "is1b", 123456789012, reflect.TypeOf(database))
	job := DatabaseJob(databasePayload)
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
			"SakuraCloud:Database",
			"SakuraCloud:Zone-is1b",
			"SakuraCloud:tag1",
			"SakuraCloud:tag2",
		})
		assert.Equal(t, para.Interfaces[0].Name, "eth0")
		assert.Equal(t, para.Interfaces[0].IPAddress, "8.8.8.8")
		assert.Equal(t, para.Interfaces[0].MacAddress, "")

	}, errFunc, errFunc)

	database.Interfaces[0].Switch.Scope = sacloud.ESCopeUser

	databasePayload = core.NewCreateHostPayload(database, "is1b", 123456789012, reflect.TypeOf(database))
	job = DatabaseJob(databasePayload)
	go job.Start(queue, option)
	jobCheckFunc(t, queue, func(v core.JobRequestAPI) {
		payload := v.GetPayload().(*core.CreateHostPayload)
		assert.NotNil(t, payload)

		para := payload.MackerelHostParam
		assert.NotNil(t, para)

		assert.Equal(t, para.Interfaces[0].Name, "eth0")
		assert.Equal(t, para.Interfaces[0].IPAddress, "192.168.0.1")
		assert.Equal(t, para.Interfaces[0].MacAddress, "")

	}, errFunc, errFunc)

}
