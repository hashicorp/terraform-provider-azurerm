package storagetaskassignments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageTaskAssignmentExecutionContext struct {
	Target  *ExecutionTarget `json:"target,omitempty"`
	Trigger ExecutionTrigger `json:"trigger"`
}
