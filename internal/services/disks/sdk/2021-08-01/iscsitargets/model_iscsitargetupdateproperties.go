package iscsitargets

type IscsiTargetUpdateProperties struct {
	Luns       *[]IscsiLun `json:"luns,omitempty"`
	StaticAcls *[]Acl      `json:"staticAcls,omitempty"`
}
