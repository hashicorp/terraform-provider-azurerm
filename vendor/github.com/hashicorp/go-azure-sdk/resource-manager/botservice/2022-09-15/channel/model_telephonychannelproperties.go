package channel

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TelephonyChannelProperties struct {
	ApiConfigurations               *[]TelephonyChannelResourceApiConfiguration `json:"apiConfigurations,omitempty"`
	CognitiveServiceRegion          *string                                     `json:"cognitiveServiceRegion,omitempty"`
	CognitiveServiceSubscriptionKey *string                                     `json:"cognitiveServiceSubscriptionKey,omitempty"`
	DefaultLocale                   *string                                     `json:"defaultLocale,omitempty"`
	IsEnabled                       *bool                                       `json:"isEnabled,omitempty"`
	PhoneNumbers                    *[]TelephonyPhoneNumbers                    `json:"phoneNumbers,omitempty"`
	PremiumSKU                      *string                                     `json:"premiumSKU,omitempty"`
}
