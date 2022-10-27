package objectreplicationpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ObjectReplicationPolicyRule struct {
	DestinationContainer string                         `json:"destinationContainer"`
	Filters              *ObjectReplicationPolicyFilter `json:"filters,omitempty"`
	RuleId               *string                        `json:"ruleId,omitempty"`
	SourceContainer      string                         `json:"sourceContainer"`
}
