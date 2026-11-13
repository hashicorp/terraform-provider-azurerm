package ddoscustompolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DdosSourcePolicyOverride struct {
	Conditions   DdosSourceMatchConditions `json:"conditions"`
	PolicyAction DdosSourcePolicyAction    `json:"policyAction"`
}
