package appserviceplans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SlowRequestsBasedTrigger struct {
	Count        *int64  `json:"count,omitempty"`
	Path         *string `json:"path,omitempty"`
	TimeInterval *string `json:"timeInterval,omitempty"`
	TimeTaken    *string `json:"timeTaken,omitempty"`
}
