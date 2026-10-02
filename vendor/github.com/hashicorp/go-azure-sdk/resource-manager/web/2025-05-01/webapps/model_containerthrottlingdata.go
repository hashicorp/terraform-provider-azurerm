package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerThrottlingData struct {
	Periods          *int64 `json:"periods,omitempty"`
	ThrottledPeriods *int64 `json:"throttledPeriods,omitempty"`
	ThrottledTime    *int64 `json:"throttledTime,omitempty"`
}
