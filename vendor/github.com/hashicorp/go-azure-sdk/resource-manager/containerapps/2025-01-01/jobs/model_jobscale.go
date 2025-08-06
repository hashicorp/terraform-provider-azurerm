package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobScale struct {
	MaxExecutions   *int64          `json:"maxExecutions,omitempty"`
	MinExecutions   *int64          `json:"minExecutions,omitempty"`
	PollingInterval *int64          `json:"pollingInterval,omitempty"`
	Rules           *[]JobScaleRule `json:"rules,omitempty"`
}
