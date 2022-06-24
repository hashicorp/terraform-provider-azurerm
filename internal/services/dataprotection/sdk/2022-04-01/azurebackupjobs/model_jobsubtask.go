package azurebackupjobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobSubTask struct {
	AdditionalDetails *map[string]string `json:"additionalDetails,omitempty"`
	TaskId            int64              `json:"taskId"`
	TaskName          string             `json:"taskName"`
	TaskProgress      *string            `json:"taskProgress,omitempty"`
	TaskStatus        string             `json:"taskStatus"`
}
