package hubs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DebugSendResult struct {
	Failure *int64                `json:"failure,omitempty"`
	Results *[]RegistrationResult `json:"results,omitempty"`
	Success *int64                `json:"success,omitempty"`
}
