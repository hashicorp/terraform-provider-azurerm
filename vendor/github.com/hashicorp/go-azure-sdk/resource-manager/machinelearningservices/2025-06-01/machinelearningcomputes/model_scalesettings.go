package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScaleSettings struct {
	MaxNodeCount                int64   `json:"maxNodeCount"`
	MinNodeCount                *int64  `json:"minNodeCount,omitempty"`
	NodeIdleTimeBeforeScaleDown *string `json:"nodeIdleTimeBeforeScaleDown,omitempty"`
}
