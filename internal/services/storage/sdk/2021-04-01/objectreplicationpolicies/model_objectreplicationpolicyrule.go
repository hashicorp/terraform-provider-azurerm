package objectreplicationpolicies

type ObjectReplicationPolicyRule struct {
	DestinationContainer string                         `json:"destinationContainer"`
	Filters              *ObjectReplicationPolicyFilter `json:"filters,omitempty"`
	RuleId               *string                        `json:"ruleId,omitempty"`
	SourceContainer      string                         `json:"sourceContainer"`
}
