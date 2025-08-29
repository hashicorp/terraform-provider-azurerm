package diagnostic

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SamplingSettings struct {
	Percentage   *float64      `json:"percentage,omitempty"`
	SamplingType *SamplingType `json:"samplingType,omitempty"`
}
