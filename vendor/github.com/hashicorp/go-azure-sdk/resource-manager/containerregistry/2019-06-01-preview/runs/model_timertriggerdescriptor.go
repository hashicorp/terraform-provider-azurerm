package runs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TimerTriggerDescriptor struct {
	ScheduleOccurrence *string `json:"scheduleOccurrence,omitempty"`
	TimerTriggerName   *string `json:"timerTriggerName,omitempty"`
}
