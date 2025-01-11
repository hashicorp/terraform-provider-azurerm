package monitors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SwitchBillingRequest struct {
	AzureResourceId *string   `json:"azureResourceId,omitempty"`
	OrganizationId  *string   `json:"organizationId,omitempty"`
	PlanData        *PlanData `json:"planData,omitempty"`
	UserEmail       string    `json:"userEmail"`
}
