package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WaitStatistics struct {
	WaitCount  *int64   `json:"waitCount,omitempty"`
	WaitTimeMs *float64 `json:"waitTimeMs,omitempty"`
	WaitType   *string  `json:"waitType,omitempty"`
}
