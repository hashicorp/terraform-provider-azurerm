package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoscaleProfile struct {
	AutoscaleType               *AutoscaleType       `json:"autoscaleType,omitempty"`
	Enabled                     bool                 `json:"enabled"`
	GracefulDecommissionTimeout *int64               `json:"gracefulDecommissionTimeout,omitempty"`
	LoadBasedConfig             *LoadBasedConfig     `json:"loadBasedConfig,omitempty"`
	ScheduleBasedConfig         *ScheduleBasedConfig `json:"scheduleBasedConfig,omitempty"`
}
