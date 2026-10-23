package containerappssessionpools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LifecycleConfiguration struct {
	CooldownPeriodInSeconds *int64         `json:"cooldownPeriodInSeconds,omitempty"`
	LifecycleType           *LifecycleType `json:"lifecycleType,omitempty"`
	MaxAlivePeriodInSeconds *int64         `json:"maxAlivePeriodInSeconds,omitempty"`
}
