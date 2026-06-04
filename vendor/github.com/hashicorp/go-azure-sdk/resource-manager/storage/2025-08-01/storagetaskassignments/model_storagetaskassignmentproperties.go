package storagetaskassignments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageTaskAssignmentProperties struct {
	Description       string                                  `json:"description"`
	Enabled           bool                                    `json:"enabled"`
	ExecutionContext  StorageTaskAssignmentExecutionContext   `json:"executionContext"`
	ProvisioningState *StorageTaskAssignmentProvisioningState `json:"provisioningState,omitempty"`
	Report            StorageTaskAssignmentReport             `json:"report"`
	RunStatus         *StorageTaskReportProperties            `json:"runStatus,omitempty"`
	TaskId            string                                  `json:"taskId"`
}
