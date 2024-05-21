package storagetaskassignments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageTaskAssignmentUpdateExecutionContext struct {
	Target  *ExecutionTargetUpdate  `json:"target,omitempty"`
	Trigger *ExecutionTriggerUpdate `json:"trigger,omitempty"`
}
