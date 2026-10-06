package storagetaskassignments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExecutionTrigger struct {
	Parameters TriggerParameters `json:"parameters"`
	Type       TriggerType       `json:"type"`
}
