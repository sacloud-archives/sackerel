package sacloud

import "time"

type Appliance struct {
	*Resource
	Class       string `json:",omitempty"`
	Name        string `json:",omitempty"`
	Description string `json:",omitempty"`
	Plan        *Resource
	//Settings
	SettingHash string `json:",omitempty"`
	//Remark      *ApplianceRemark `json:",omitempty"`
	*EAvailability
	Instance     *EServerInstanceStatus `json:",omitempty"`
	ServiceClass string                 `json:",omitempty"`
	CreatedAt    *time.Time             `json:",omitempty"`
	Icon         *Icon                  `json:",omitempty"`
	Switch       *Switch                `json:",omitempty"`
	Interfaces   []Interface            `json:",omitempty"`
	*TagsType
}

//HACK Appliance:Zone.IDがRoute/LoadBalancerの場合でデータ型が異なるため
//それぞれのstruct定義でZoneだけ上書きした構造体を定義して使う

type ApplianceRemarkBase struct {
	Servers []interface{}
	Switch  *ApplianceRemarkSwitch `json:",omitempty"`
	//Zone *Resource `json:",omitempty"`
	VRRP    *ApplianceRemarkVRRP    `json:",omitempty"`
	Network *ApplianceRemarkNetwork `json:",omitempty"`
	//Plan    *Resource
}

//type ApplianceServer struct {
//	IPAddress string `json:",omitempty"`
//}

type ApplianceRemarkSwitch struct {
	ID    string `json:",omitempty"`
	Scope string `json:",omitempty"`
}

type ApplianceRemarkVRRP struct {
	VRID int `json:",omitempty"`
}

type ApplianceRemarkNetwork struct {
	NetworkMaskLen int    `json:",omitempty"`
	DefaultRoute   string `json:",omitempty"`
}
