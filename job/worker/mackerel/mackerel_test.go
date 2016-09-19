package mackerel

import (
	mkr "github.com/mackerelio/mackerel-client-go"
	"github.com/sacloud/sackerel/job/core"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

var (
	testSetupHandlers    []func()
	testTearDownHandlers []func()
	mockOption           *core.Option
	mockBlockQueue       *core.Queue
	mackerelClient       *mkr.Client
)

func TestMain(m *testing.M) {
	//環境変数にトークン/シークレットがある場合のみテスト実施
	apiKey := os.Getenv("MACKEREL_APIKEY")

	if apiKey == "" {
		log.Println("Please Set ENV 'MACKEREL_APIKEY'")
		os.Exit(0) // exit normal
	}

	// setup mock option
	mockOption = core.NewOption()
	mockOption.MackerelOption.APIKey = apiKey

	// setup mock queue
	mockBlockQueue = core.NewQueue(1, 1, 1, 1)

	mackerelClient = getClient(mockOption)

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

func TestGetOrganization(t *testing.T) {

	org, err := mackerelClient.GetOrg()

	assert.NoError(t, err)
	assert.NotNil(t, org)
	assert.NotEmpty(t, org.Name)

}
