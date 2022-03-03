package objectreplicationpolicies

type ObjectReplicationPolicy struct {
	Id         *string                            `json:"id,omitempty"`
	Name       *string                            `json:"name,omitempty"`
	Properties *ObjectReplicationPolicyProperties `json:"properties,omitempty"`
	Type       *string                            `json:"type,omitempty"`
}
