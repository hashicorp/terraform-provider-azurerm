package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobConstraints struct {
	MaxTaskRetryCount *int64  `json:"maxTaskRetryCount,omitempty"`
	MaxWallClockTime  *string `json:"maxWallClockTime,omitempty"`
}
