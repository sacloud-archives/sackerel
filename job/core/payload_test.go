package core

import (
	mkr "github.com/mackerelio/mackerel-client-go"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestCreateHostPayload_GetFindParam(t *testing.T) {
	s := NewCreateHostPayload("", "tk1v", 123456789012, reflect.TypeOf(""))
	p := s.GetFindParam()
	mackerelName := s.GenerateMackerelName()

	assert.NotNil(t, p)

	// セットするもの/しないもの
	assert.NotEmpty(t, p.CustomIdentifier)
	assert.NotEmpty(t, p.Statuses)
	assert.NotEmpty(t, p.Name)
	assert.Empty(t, p.Roles)
	assert.Empty(t, p.Service)

	// セットする値
	assert.Equal(t, p.Name, mackerelName)
	assert.Equal(t, p.CustomIdentifier, mackerelName)
	assert.Len(t, p.Statuses, 4) // mackerel側で取り得るステータス全て
}

func TestCreateHostPayload_IsRoleUpdated(t *testing.T) {
	s := NewCreateHostPayload("", "tk1v", 123456789012, reflect.TypeOf(""))

	// Mackerelホスト(比較元)が設定されていない場合はエラーとする
	res, err := s.IsRoleUpdated()
	assert.False(t, res)
	assert.Error(t, err)
	err = nil

	// Mackerelホスト
	s.MackerelHost = &mkr.Host{
		ID: "xxxxxxxxxxx",
	}

	// Createパラメータ(比較先)が設定されていない場合はエラーとする
	res, err = s.IsRoleUpdated()
	assert.False(t, res)
	assert.Error(t, err)
	err = nil

	// Createパラメータ
	s.MackerelHostParam = &mkr.CreateHostParam{}

	// 両方ともロールなしの場合は更新なし判定
	res, err = s.IsRoleUpdated()
	assert.False(t, res)
	assert.NoError(t, err)

	// 双方のロール数が違う場合は更新あり
	s.MackerelHost.Roles = map[string][]string{
		"SakuraCloud": {"tag1", "tag2"},
	}
	s.MackerelHostParam.RoleFullnames = []string{"SakuraCloud:tag1"}

	res, err = s.IsRoleUpdated()
	assert.True(t, res)
	assert.NoError(t, err)

	// 双方のロール数が同じ場合は内容で判定
	s.MackerelHost.Roles = map[string][]string{
		"SakuraCloud": {"tag1", "tag2"},
	}
	s.MackerelHostParam.RoleFullnames = []string{"SakuraCloud:tag1", "SakuraCloud:tag2"}

	res, err = s.IsRoleUpdated()
	assert.False(t, res)
	assert.NoError(t, err)

	// 順不同
	s.MackerelHost.Roles = map[string][]string{
		"SakuraCloud": {"tag1", "tag2"},
	}
	s.MackerelHostParam.RoleFullnames = []string{"SakuraCloud:tag2", "SakuraCloud:tag1"}

	res, err = s.IsRoleUpdated()
	assert.False(t, res)
	assert.NoError(t, err)

	s.MackerelHost.Roles = map[string][]string{
		"SakuraCloud": {"tag1", "tag2"},
	}
	s.MackerelHostParam.RoleFullnames = []string{"SakuraCloud:tag1", "SakuraCloud:tag2aaaaa"}

	res, err = s.IsRoleUpdated()
	assert.True(t, res)
	assert.NoError(t, err)

}

func TestCreateHostPayload_IsStatusUpdated(t *testing.T) {
	s := NewCreateHostPayload("", "tk1v", 123456789012, reflect.TypeOf(""))

	// Mackerelホスト(比較元)が設定されていない場合はエラーとする
	res, err := s.IsStatusUpdated()
	assert.False(t, res)
	assert.Error(t, err)
	err = nil

	// Mackerelホスト
	s.MackerelHost = &mkr.Host{
		ID:     "xxxxxxxxxxx",
		Status: string(MackerelHostStatusWorking),
	}

	// ホストステータス(比較先)が設定されていない場合はエラーとする
	s.MackerelHostStatus = MackerelHostStatus("")
	res, err = s.IsStatusUpdated()
	assert.False(t, res)
	assert.Error(t, err)
	err = nil

	s.MackerelHostStatus = MackerelHostStatusWorking

	// 値比較
	res, err = s.IsStatusUpdated()
	assert.False(t, res)
	assert.NoError(t, err)
	err = nil

	s.MackerelHostStatus = MackerelHostStatusMaintenance
	res, err = s.IsStatusUpdated()
	assert.True(t, res)
	assert.NoError(t, err)
	err = nil

}
