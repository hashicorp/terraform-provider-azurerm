package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Schedule struct {
	Count     int64         `json:"count"`
	Days      []ScheduleDay `json:"days"`
	EndTime   string        `json:"endTime"`
	StartTime string        `json:"startTime"`
}
