package retentionpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RetentionPolicyProperties struct {
	ProvisioningState *ProvisioningState        `json:"provisioningState,omitempty"`
	RetentionPolicies *[]RetentionPolicyDetails `json:"retentionPolicies,omitempty"`
}
