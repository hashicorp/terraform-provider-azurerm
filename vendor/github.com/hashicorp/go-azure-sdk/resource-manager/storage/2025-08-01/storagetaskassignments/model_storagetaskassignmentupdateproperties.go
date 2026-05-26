package storagetaskassignments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageTaskAssignmentUpdateProperties struct {
	Description       *string                                      `json:"description,omitempty"`
	Enabled           *bool                                        `json:"enabled,omitempty"`
	ExecutionContext  *StorageTaskAssignmentUpdateExecutionContext `json:"executionContext,omitempty"`
	ProvisioningState *StorageTaskAssignmentProvisioningState      `json:"provisioningState,omitempty"`
	Report            *StorageTaskAssignmentUpdateReport           `json:"report,omitempty"`
	RunStatus         *StorageTaskReportProperties                 `json:"runStatus,omitempty"`
	TaskId            *string                                      `json:"taskId,omitempty"`
}
