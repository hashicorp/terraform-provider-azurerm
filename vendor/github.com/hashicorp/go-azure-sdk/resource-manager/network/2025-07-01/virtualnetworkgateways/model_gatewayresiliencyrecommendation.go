package virtualnetworkgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GatewayResiliencyRecommendation struct {
	CallToActionLink    *string `json:"callToActionLink,omitempty"`
	CallToActionText    *string `json:"callToActionText,omitempty"`
	RecommendationId    *string `json:"recommendationId,omitempty"`
	RecommendationText  *string `json:"recommendationText,omitempty"`
	RecommendationTitle *string `json:"recommendationTitle,omitempty"`
	Severity            *string `json:"severity,omitempty"`
}
