package channel

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TelephonyChannelResourceApiConfiguration struct {
	CognitiveServiceRegion          *string `json:"cognitiveServiceRegion,omitempty"`
	CognitiveServiceResourceId      *string `json:"cognitiveServiceResourceId,omitempty"`
	CognitiveServiceSubscriptionKey *string `json:"cognitiveServiceSubscriptionKey,omitempty"`
	DefaultLocale                   *string `json:"defaultLocale,omitempty"`
	Id                              *string `json:"id,omitempty"`
	ProviderName                    *string `json:"providerName,omitempty"`
}
