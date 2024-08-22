package monitors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NewRelicSingleSignOnProperties struct {
	EnterpriseAppId   *string             `json:"enterpriseAppId,omitempty"`
	ProvisioningState *ProvisioningState  `json:"provisioningState,omitempty"`
	SingleSignOnState *SingleSignOnStates `json:"singleSignOnState,omitempty"`
	SingleSignOnUrl   *string             `json:"singleSignOnUrl,omitempty"`
}
