package onlinedeployment

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProbeSettings struct {
	FailureThreshold *int64  `json:"failureThreshold,omitempty"`
	InitialDelay     *string `json:"initialDelay,omitempty"`
	Period           *string `json:"period,omitempty"`
	SuccessThreshold *int64  `json:"successThreshold,omitempty"`
	Timeout          *string `json:"timeout,omitempty"`
}
