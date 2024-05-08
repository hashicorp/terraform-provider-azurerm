package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LoadBasedConfig struct {
	CooldownPeriod *int64        `json:"cooldownPeriod,omitempty"`
	MaxNodes       int64         `json:"maxNodes"`
	MinNodes       int64         `json:"minNodes"`
	PollInterval   *int64        `json:"pollInterval,omitempty"`
	ScalingRules   []ScalingRule `json:"scalingRules"`
}
