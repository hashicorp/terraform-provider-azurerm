package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobProperties struct {
	Description *string      `json:"description,omitempty"`
	Schedule    *JobSchedule `json:"schedule,omitempty"`
	Version     *int64       `json:"version,omitempty"`
}
