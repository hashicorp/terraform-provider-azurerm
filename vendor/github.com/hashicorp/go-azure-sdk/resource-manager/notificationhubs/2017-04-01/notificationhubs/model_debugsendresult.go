package notificationhubs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DebugSendResult struct {
	Failure *float64     `json:"failure,omitempty"`
	Results *interface{} `json:"results,omitempty"`
	Success *float64     `json:"success,omitempty"`
}
