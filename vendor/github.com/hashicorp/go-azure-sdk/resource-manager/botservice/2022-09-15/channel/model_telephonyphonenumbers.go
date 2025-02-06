package channel

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TelephonyPhoneNumbers struct {
	AcsEndpoint                     *string `json:"acsEndpoint,omitempty"`
	AcsResourceId                   *string `json:"acsResourceId,omitempty"`
	AcsSecret                       *string `json:"acsSecret,omitempty"`
	CognitiveServiceRegion          *string `json:"cognitiveServiceRegion,omitempty"`
	CognitiveServiceResourceId      *string `json:"cognitiveServiceResourceId,omitempty"`
	CognitiveServiceSubscriptionKey *string `json:"cognitiveServiceSubscriptionKey,omitempty"`
	DefaultLocale                   *string `json:"defaultLocale,omitempty"`
	Id                              *string `json:"id,omitempty"`
	OfferType                       *string `json:"offerType,omitempty"`
	PhoneNumber                     *string `json:"phoneNumber,omitempty"`
}
