package backuprestores

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SelectiveKeyRestoreOperation struct {
	EndTime       *int64           `json:"endTime,omitempty"`
	Error         *Error           `json:"error,omitempty"`
	JobId         *string          `json:"jobId,omitempty"`
	StartTime     *int64           `json:"startTime,omitempty"`
	Status        *OperationStatus `json:"status,omitempty"`
	StatusDetails *string          `json:"statusDetails,omitempty"`
}
