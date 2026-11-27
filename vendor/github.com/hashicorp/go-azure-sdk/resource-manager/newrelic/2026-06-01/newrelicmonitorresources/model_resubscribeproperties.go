package newrelicmonitorresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResubscribeProperties struct {
	OfferId        *string `json:"offerId,omitempty"`
	OrganizationId *string `json:"organizationId,omitempty"`
	PlanId         *string `json:"planId,omitempty"`
	PublisherId    *string `json:"publisherId,omitempty"`
	ResourceGroup  *string `json:"resourceGroup,omitempty"`
	SubscriptionId *string `json:"subscriptionId,omitempty"`
	TermId         *string `json:"termId,omitempty"`
}
