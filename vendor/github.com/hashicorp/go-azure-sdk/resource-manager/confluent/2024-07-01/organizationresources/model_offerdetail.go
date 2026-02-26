package organizationresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OfferDetail struct {
	Id              string           `json:"id"`
	PlanId          string           `json:"planId"`
	PlanName        string           `json:"planName"`
	PrivateOfferId  *string          `json:"privateOfferId,omitempty"`
	PrivateOfferIds *[]string        `json:"privateOfferIds,omitempty"`
	PublisherId     string           `json:"publisherId"`
	Status          *SaaSOfferStatus `json:"status,omitempty"`
	TermId          *string          `json:"termId,omitempty"`
	TermUnit        string           `json:"termUnit"`
}
