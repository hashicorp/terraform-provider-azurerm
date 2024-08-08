package deployments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ThrottlingRule struct {
	Count                    *float64               `json:"count,omitempty"`
	DynamicThrottlingEnabled *bool                  `json:"dynamicThrottlingEnabled,omitempty"`
	Key                      *string                `json:"key,omitempty"`
	MatchPatterns            *[]RequestMatchPattern `json:"matchPatterns,omitempty"`
	MinCount                 *float64               `json:"minCount,omitempty"`
	RenewalPeriod            *float64               `json:"renewalPeriod,omitempty"`
}
