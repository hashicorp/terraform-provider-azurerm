package workflowtriggers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkflowTriggerReference struct {
	FlowName    *string `json:"flowName,omitempty"`
	Id          *string `json:"id,omitempty"`
	Name        *string `json:"name,omitempty"`
	TriggerName *string `json:"triggerName,omitempty"`
	Type        *string `json:"type,omitempty"`
}
