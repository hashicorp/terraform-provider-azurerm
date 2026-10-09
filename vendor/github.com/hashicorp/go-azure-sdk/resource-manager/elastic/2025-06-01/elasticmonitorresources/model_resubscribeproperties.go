package elasticmonitorresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResubscribeProperties struct {
	OrganizationId *string `json:"organizationId,omitempty"`
	PlanId         *string `json:"planId,omitempty"`
	ResourceGroup  *string `json:"resourceGroup,omitempty"`
	SubscriptionId *string `json:"subscriptionId,omitempty"`
	Term           *string `json:"term,omitempty"`
}
