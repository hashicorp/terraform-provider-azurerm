package organizations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LiftrBaseSingleSignOnProperties struct {
	AadDomains        *[]string                  `json:"aadDomains,omitempty"`
	EnterpriseAppId   *string                    `json:"enterpriseAppId,omitempty"`
	ProvisioningState *ResourceProvisioningState `json:"provisioningState,omitempty"`
	SingleSignOnState *SingleSignOnStates        `json:"singleSignOnState,omitempty"`
	SingleSignOnURL   *string                    `json:"singleSignOnUrl,omitempty"`
}
