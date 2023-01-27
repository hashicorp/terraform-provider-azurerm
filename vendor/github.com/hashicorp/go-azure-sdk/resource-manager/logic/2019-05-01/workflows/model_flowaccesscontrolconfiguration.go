package workflows

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FlowAccessControlConfiguration struct {
	Actions            *FlowAccessControlConfigurationPolicy `json:"actions,omitempty"`
	Contents           *FlowAccessControlConfigurationPolicy `json:"contents,omitempty"`
	Triggers           *FlowAccessControlConfigurationPolicy `json:"triggers,omitempty"`
	WorkflowManagement *FlowAccessControlConfigurationPolicy `json:"workflowManagement,omitempty"`
}
