package sacloud

import (
	"fmt"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/sackerel/job/core"
	"reflect"
)

func doActionPerZone(option *core.Option, sacloudAPIFunc func(*api.Client) error) error {
	for _, zone := range option.SakuraCloudOption.Zone {

		client := getClient(option, zone)
		// call API func per zone.
		err := sacloudAPIFunc(client)
		if err != nil {
			return err
		}
	}

	return nil
}

func getClient(option *core.Option, zone string) *api.Client {
	o := option.SakuraCloudOption

	client := api.NewClient(o.AccessToken, o.AccessTokenSecret, zone)
	client.TraceMode = o.TraceMode

	return client

}

func hasIgnoreTagWithInfoLog(target interface{}, ignoreTag string, queue *core.Queue, id int64, name string, t reflect.Type) bool {

	if tagsType, ok := target.(tagsHolder); ok {
		if tagsType.HasTag(ignoreTag) {
			targetName := fmt.Sprintf("(%s)", t.Name())
			queue.PushInfo(
				fmt.Sprintf(
					"Ignore target  %-15s => SakuraID:[%d] / Name:[%s]",
					targetName,
					id,
					name,
				),
			)
			return true
		}
	}

	return false
}

type tagsHolder interface {
	HasTag(string) bool
}
