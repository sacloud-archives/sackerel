package core

import (
	"fmt"
	mkr "github.com/mackerelio/mackerel-client-go"
	"reflect"
)

//-----------------------------------------------------------------------------
// Payload common data types
//-----------------------------------------------------------------------------

// MackerelHostStatus Mackerel上でのホストステータスを表す
type MackerelHostStatus string

var (
	// MackerelHostStatusStandby "standby"ホストステータス
	MackerelHostStatusStandby = MackerelHostStatus("standby")

	// MackerelHostStatusWorking "working"ホストステータス
	MackerelHostStatusWorking = MackerelHostStatus("working")

	// MackerelHostStatusMaintenance "maintenance"ホストステータス
	MackerelHostStatusMaintenance = MackerelHostStatus("maintenance")

	// MackerelHostStatusPowerOff "poweroff"ホストステータス
	MackerelHostStatusPowerOff = MackerelHostStatus("poweroff")
)

//-----------------------------------------------------------------------------
// Payload interfaces
//-----------------------------------------------------------------------------

// MackerelFindParamHolder Mackerelからのホスト検索用パラメータを保持していることを表すインターフェース
type MackerelFindParamHolder interface {
	GetFindParam() *mkr.FindHostsParam
}

// SourcePayloadHolder SourcePayloadを保持していることを表すインターフェース
type SourcePayloadHolder interface {
	GetSourcePayload() *SourcePayload
}

//-----------------------------------------------------------------------------
// SourceMetaPayload
//-----------------------------------------------------------------------------

// SourcePayload 連携元データを内包するペイロード
type SourcePayload struct {
	SacloudSource      interface{}
	SacloudZone        string
	SacloudResourceID  int64
	SourceType         reflect.Type
	MackerelID         string
	MackerelHost       *mkr.Host
	MackerelHostStatus MackerelHostStatus
}

// GetSourcePayload ペイロードの取得
func (p *SourcePayload) GetSourcePayload() *SourcePayload {
	return p
}

// GenerateMackerelName さくらのクラウド上のリソース定義を元にMackerel上でのホスト名を生成する
func (p *SourcePayload) GenerateMackerelName() string {
	resourceID := p.SacloudResourceID
	if resourceID <= 0 {
		return ""
	}
	return fmt.Sprintf("sakuracloud-%s-%d", p.SacloudZone, resourceID)
}

// NewSourcePayload SourcePayloadの新規作成
func NewSourcePayload(source interface{}, zone string, resourceID int64, sourceType reflect.Type) *SourcePayload {
	return &SourcePayload{
		SacloudSource:      source,
		SacloudZone:        zone,
		SacloudResourceID:  resourceID,
		SourceType:         sourceType,
		MackerelHostStatus: MackerelHostStatusWorking,
	}
}
