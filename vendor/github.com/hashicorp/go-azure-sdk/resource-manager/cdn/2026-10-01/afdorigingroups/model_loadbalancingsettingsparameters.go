package afdorigingroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LoadBalancingSettingsParameters struct {
	AdditionalLatencyInMilliseconds *int64 `json:"additionalLatencyInMilliseconds,omitempty"`
	SampleSize                      *int64 `json:"sampleSize,omitempty"`
	SuccessfulSamplesRequired       *int64 `json:"successfulSamplesRequired,omitempty"`
}
