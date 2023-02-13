package frontdoors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LoadBalancingSettingsProperties struct {
	AdditionalLatencyMilliseconds *int64                  `json:"additionalLatencyMilliseconds,omitempty"`
	ResourceState                 *FrontDoorResourceState `json:"resourceState,omitempty"`
	SampleSize                    *int64                  `json:"sampleSize,omitempty"`
	SuccessfulSamplesRequired     *int64                  `json:"successfulSamplesRequired,omitempty"`
}
