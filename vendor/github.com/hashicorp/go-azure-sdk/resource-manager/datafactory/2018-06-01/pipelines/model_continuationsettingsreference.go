package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContinuationSettingsReference struct {
	ContinuationTtlInMinutes *interface{} `json:"continuationTtlInMinutes,omitempty"`
	CustomizedCheckpointKey  *interface{} `json:"customizedCheckpointKey,omitempty"`
	IdleCondition            *interface{} `json:"idleCondition,omitempty"`
}
