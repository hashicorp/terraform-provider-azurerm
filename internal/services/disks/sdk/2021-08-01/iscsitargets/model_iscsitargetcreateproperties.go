package iscsitargets

type IscsiTargetCreateProperties struct {
	AclMode    IscsiTargetAclMode `json:"aclMode"`
	Luns       *[]IscsiLun        `json:"luns,omitempty"`
	StaticAcls *[]Acl             `json:"staticAcls,omitempty"`
	TargetIqn  *string            `json:"targetIqn,omitempty"`
}
