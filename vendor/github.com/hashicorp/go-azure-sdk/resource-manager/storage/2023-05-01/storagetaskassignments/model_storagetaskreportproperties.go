package storagetaskassignments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageTaskReportProperties struct {
	FinishTime             *string        `json:"finishTime,omitempty"`
	ObjectFailedCount      *string        `json:"objectFailedCount,omitempty"`
	ObjectsOperatedOnCount *string        `json:"objectsOperatedOnCount,omitempty"`
	ObjectsSucceededCount  *string        `json:"objectsSucceededCount,omitempty"`
	ObjectsTargetedCount   *string        `json:"objectsTargetedCount,omitempty"`
	RunResult              *RunResult     `json:"runResult,omitempty"`
	RunStatusEnum          *RunStatusEnum `json:"runStatusEnum,omitempty"`
	RunStatusError         *string        `json:"runStatusError,omitempty"`
	StartTime              *string        `json:"startTime,omitempty"`
	StorageAccountId       *string        `json:"storageAccountId,omitempty"`
	SummaryReportPath      *string        `json:"summaryReportPath,omitempty"`
	TaskAssignmentId       *string        `json:"taskAssignmentId,omitempty"`
	TaskId                 *string        `json:"taskId,omitempty"`
	TaskVersion            *string        `json:"taskVersion,omitempty"`
}
