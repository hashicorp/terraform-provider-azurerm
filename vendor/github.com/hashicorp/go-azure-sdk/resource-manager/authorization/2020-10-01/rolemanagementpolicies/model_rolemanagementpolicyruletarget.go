package rolemanagementpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoleManagementPolicyRuleTarget struct {
	Caller              *string   `json:"caller,omitempty"`
	EnforcedSettings    *[]string `json:"enforcedSettings,omitempty"`
	InheritableSettings *[]string `json:"inheritableSettings,omitempty"`
	Level               *string   `json:"level,omitempty"`
	Operations          *[]string `json:"operations,omitempty"`
	TargetObjects       *[]string `json:"targetObjects,omitempty"`
}
