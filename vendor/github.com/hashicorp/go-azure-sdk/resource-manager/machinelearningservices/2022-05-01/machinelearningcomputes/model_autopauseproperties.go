package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoPauseProperties struct {
	DelayInMinutes *int64 `json:"delayInMinutes,omitempty"`
	Enabled        *bool  `json:"enabled,omitempty"`
}
