package organizations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LiftrBaseOfferDetails struct {
	OfferId     string  `json:"offerId"`
	PlanId      string  `json:"planId"`
	PlanName    *string `json:"planName,omitempty"`
	PublisherId string  `json:"publisherId"`
	TermId      *string `json:"termId,omitempty"`
	TermUnit    *string `json:"termUnit,omitempty"`
}
