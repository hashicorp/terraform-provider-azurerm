package indexes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MagnitudeScoringParameters struct {
	BoostingRangeEnd         float64 `json:"boostingRangeEnd"`
	BoostingRangeStart       float64 `json:"boostingRangeStart"`
	ConstantBoostBeyondRange *bool   `json:"constantBoostBeyondRange,omitempty"`
}
