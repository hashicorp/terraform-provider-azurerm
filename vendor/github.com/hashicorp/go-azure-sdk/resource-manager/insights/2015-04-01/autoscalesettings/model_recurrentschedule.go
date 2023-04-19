package autoscalesettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RecurrentSchedule struct {
	Days     []string `json:"days"`
	Hours    []int64  `json:"hours"`
	Minutes  []int64  `json:"minutes"`
	TimeZone string   `json:"timeZone"`
}
