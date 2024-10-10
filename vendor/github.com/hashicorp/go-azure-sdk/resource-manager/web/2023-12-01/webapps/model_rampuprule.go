package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RampUpRule struct {
	ActionHostName            *string  `json:"actionHostName,omitempty"`
	ChangeDecisionCallbackURL *string  `json:"changeDecisionCallbackUrl,omitempty"`
	ChangeIntervalInMinutes   *int64   `json:"changeIntervalInMinutes,omitempty"`
	ChangeStep                *float64 `json:"changeStep,omitempty"`
	MaxReroutePercentage      *float64 `json:"maxReroutePercentage,omitempty"`
	MinReroutePercentage      *float64 `json:"minReroutePercentage,omitempty"`
	Name                      *string  `json:"name,omitempty"`
	ReroutePercentage         *float64 `json:"reroutePercentage,omitempty"`
}
