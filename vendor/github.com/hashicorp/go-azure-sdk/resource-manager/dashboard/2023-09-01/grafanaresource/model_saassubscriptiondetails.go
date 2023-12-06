package grafanaresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SaasSubscriptionDetails struct {
	OfferId     *string           `json:"offerId,omitempty"`
	PlanId      *string           `json:"planId,omitempty"`
	PublisherId *string           `json:"publisherId,omitempty"`
	Term        *SubscriptionTerm `json:"term,omitempty"`
}
