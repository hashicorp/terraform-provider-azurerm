package elasticmonitorresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PlanDetails struct {
	OfferID     *string `json:"offerID,omitempty"`
	PlanID      *string `json:"planID,omitempty"`
	PlanName    *string `json:"planName,omitempty"`
	PublisherID *string `json:"publisherID,omitempty"`
	TermID      *string `json:"termID,omitempty"`
}
