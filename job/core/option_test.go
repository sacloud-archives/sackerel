package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateOption(t *testing.T) {

	opt := NewOption()

	errors := opt.Validate()
	assert.True(t, len(errors) > 0)

	// AccessToken
	opt.SakuraCloudOption.AccessToken = "test"

	errors = opt.Validate()
	assert.True(t, len(errors) > 0)

	// AccessTokenSecret
	opt.SakuraCloudOption.AccessTokenSecret = "test"

	errors = opt.Validate()
	assert.True(t, len(errors) > 0)

	// Mackerel APIKey
	opt.MackerelOption.APIKey = "test"

	errors = opt.Validate()
	assert.Empty(t, errors)

}
