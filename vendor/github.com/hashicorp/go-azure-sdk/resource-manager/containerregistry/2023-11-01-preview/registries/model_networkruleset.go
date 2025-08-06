package registries

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkRuleSet struct {
	DefaultAction DefaultAction `json:"defaultAction"`
	IPRules       *[]IPRule     `json:"ipRules,omitempty"`
}
