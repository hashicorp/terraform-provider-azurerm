package workflowrunactions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RunActionCorrelation struct {
	ActionTrackingId *string   `json:"actionTrackingId,omitempty"`
	ClientKeywords   *[]string `json:"clientKeywords,omitempty"`
	ClientTrackingId *string   `json:"clientTrackingId,omitempty"`
}
