package objectreplicationpolicies

type ObjectReplicationPolicyFilter struct {
	MinCreationTime *string   `json:"minCreationTime,omitempty"`
	PrefixMatch     *[]string `json:"prefixMatch,omitempty"`
}
