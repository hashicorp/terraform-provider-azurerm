package virtualnetworkgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResiliencyRecommendationComponents struct {
	CurrentScore    *string                            `json:"currentScore,omitempty"`
	MaxScore        *string                            `json:"maxScore,omitempty"`
	Name            *string                            `json:"name,omitempty"`
	Recommendations *[]GatewayResiliencyRecommendation `json:"recommendations,omitempty"`
}
