package schedule

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TableVerticalLimitSettings struct {
	EnableEarlyTermination *bool    `json:"enableEarlyTermination,omitempty"`
	ExitScore              *float64 `json:"exitScore,omitempty"`
	MaxConcurrentTrials    *int64   `json:"maxConcurrentTrials,omitempty"`
	MaxCoresPerTrial       *int64   `json:"maxCoresPerTrial,omitempty"`
	MaxTrials              *int64   `json:"maxTrials,omitempty"`
	Timeout                *string  `json:"timeout,omitempty"`
	TrialTimeout           *string  `json:"trialTimeout,omitempty"`
}
