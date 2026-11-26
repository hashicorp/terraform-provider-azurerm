package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoScaleProperties struct {
	Enabled      *bool  `json:"enabled,omitempty"`
	MaxNodeCount *int64 `json:"maxNodeCount,omitempty"`
	MinNodeCount *int64 `json:"minNodeCount,omitempty"`
}
