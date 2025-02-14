package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContinuationSettingsReference struct {
	ContinuationTtlInMinutes *int64  `json:"continuationTtlInMinutes,omitempty"`
	CustomizedCheckpointKey  *string `json:"customizedCheckpointKey,omitempty"`
	IdleCondition            *string `json:"idleCondition,omitempty"`
}
