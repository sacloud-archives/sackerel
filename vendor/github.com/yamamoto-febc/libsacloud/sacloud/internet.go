package sacloud

import (
	"time"
)

type Internet struct {
	*Resource
	Index          int        `json:",omitempty"`
	Name           string     `json:",omitempty"`
	Description    string     `json:",omitempty"`
	BandWidthMbps  int        `json:",omitempty"`
	NetworkMaskLen int        `json:",omitempty"`
	Scope          EScope     `json:",omitempty"`
	ServiceClass   string     `json:",omitempty"`
	CreatedAt      *time.Time `json:",omitempty"`
	Icon           *Icon      `json:",omitempty"`

	Switch *Switch `json:",omitempty"`
	*TagsType

	//TODO Zone
	// Zone           *Zone      `json:",omitempty"`

}
