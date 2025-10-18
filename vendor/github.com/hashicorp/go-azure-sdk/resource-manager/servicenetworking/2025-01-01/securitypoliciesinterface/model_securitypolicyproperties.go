package securitypoliciesinterface

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityPolicyProperties struct {
	PolicyType        *PolicyType        `json:"policyType,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	WafPolicy         *WafPolicy         `json:"wafPolicy,omitempty"`
}
