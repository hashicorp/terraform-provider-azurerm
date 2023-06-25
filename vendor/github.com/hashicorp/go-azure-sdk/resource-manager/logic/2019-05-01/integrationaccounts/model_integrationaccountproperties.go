package integrationaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationAccountProperties struct {
	IntegrationServiceEnvironment *ResourceReference `json:"integrationServiceEnvironment,omitempty"`
	State                         *WorkflowState     `json:"state,omitempty"`
}
