package allpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AllPoliciesContractProperties struct {
	ComplianceState   *PolicyComplianceState `json:"complianceState,omitempty"`
	ReferencePolicyId *string                `json:"referencePolicyId,omitempty"`
}
