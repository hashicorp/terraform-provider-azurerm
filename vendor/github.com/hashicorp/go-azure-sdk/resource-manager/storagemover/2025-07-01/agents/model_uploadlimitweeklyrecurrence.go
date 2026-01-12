package agents

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UploadLimitWeeklyRecurrence struct {
	Days        []DayOfWeek `json:"days"`
	EndTime     Time        `json:"endTime"`
	LimitInMbps int64       `json:"limitInMbps"`
	StartTime   Time        `json:"startTime"`
}
