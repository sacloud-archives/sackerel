package exchange

import (
	"github.com/sacloud/libsacloud/sacloud"
	"github.com/sacloud/sackerel/job/core"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

var loadbalancer = &sacloud.LoadBalancer{
	Appliance: &sacloud.Appliance{
		Resource: &sacloud.Resource{ID: 123456789012},
		Name:     "Test",
		TagsType: &sacloud.TagsType{
			Tags: []string{"tag1", "tag2"},
		},
	},
	Remark: &sacloud.LoadBalancerRemark{
		ApplianceRemarkBase: &sacloud.ApplianceRemarkBase{
			Servers: []interface{}{
				map[string]interface{}{
					"IPAddress": "192.168.0.1",
				},
				map[string]interface{}{
					"IPAddress": "192.168.0.2",
				},
			},
		},
	},
}

func TestLoadBalancerJob(t *testing.T) {
	if !baseExchangeJobTest(t) {
		return
	}

	lbPayload := core.NewCreateHostPayload(loadbalancer, "is1b", 123456789012, reflect.TypeOf(loadbalancer))
	job := LoadBalancerJob(lbPayload)
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
			"SakuraCloud:LoadBalancer",
			"SakuraCloud:Zone-is1b",
			"SakuraCloud:tag1",
			"SakuraCloud:tag2",
		})
		assert.Equal(t, para.Interfaces[0].Name, "eth0")
		assert.Equal(t, para.Interfaces[0].IPAddress, "192.168.0.1")
		assert.Equal(t, para.Interfaces[0].MacAddress, "")

		assert.Equal(t, para.Interfaces[1].Name, "eth1")
		assert.Equal(t, para.Interfaces[1].IPAddress, "192.168.0.2")
		assert.Equal(t, para.Interfaces[1].MacAddress, "")

	}, errFunc, errFunc)

}
