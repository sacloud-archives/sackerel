package sacloud

import (
	"github.com/sacloud/sackerel/job/core"
	"log"
	"os"
	"testing"
)

var testSetupHandlers []func()
var testTearDownHandlers []func()
var mockOption *core.Option
var mockBlockQueue *core.Queue

func TestMain(m *testing.M) {
	//環境変数にトークン/シークレットがある場合のみテスト実施
	accessToken := os.Getenv("SAKURACLOUD_ACCESS_TOKEN")
	accessTokenSecret := os.Getenv("SAKURACLOUD_ACCESS_TOKEN_SECRET")

	if accessToken == "" || accessTokenSecret == "" {
		log.Println("Please Set ENV 'SAKURACLOUD_ACCESS_TOKEN' and 'SAKURACLOUD_ACCESS_TOKEN_SECRET'")
		os.Exit(0) // exit normal
	}

	// setup mock option
	mockOption = core.NewOption()
	mockOption.SakuraCloudOption.AccessToken = accessToken
	mockOption.SakuraCloudOption.AccessTokenSecret = accessTokenSecret

	// setup mock queue
	mockBlockQueue = core.NewQueue(1, 1, 1, 1)

	// setup test
	for _, f := range testSetupHandlers {
		f()
	}

	ret := m.Run()

	// teardown
	for _, f := range testTearDownHandlers {
		f()
	}

	os.Exit(ret)
}

func TestDetectServer(t *testing.T) {

	//work := NewDetectServerWork(mockOption.SakuraCloudOption).(*DetectServerWork)
	//
	//assert.NotNil(t, work)
	//
	//go work.detect(mockBlockQueue)
	//
	//select {
	//case w := <-mockBlockQueue.WorkRequest:
	//	name := w.GetName()
	//	assert.Equal(t, name, "found-server")
	//	//payload := w.GetPayload()
	//	//assert.NotNil(t, payload)
	//
	//case <-time.After(1 * time.Minute):
	//	assert.Fail(t, "Timeout [%s]", "detect")
	//}

}
