package jobexecutions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobExecutionTarget struct {
	DatabaseName *string        `json:"databaseName,omitempty"`
	ServerName   *string        `json:"serverName,omitempty"`
	Type         *JobTargetType `json:"type,omitempty"`
}
