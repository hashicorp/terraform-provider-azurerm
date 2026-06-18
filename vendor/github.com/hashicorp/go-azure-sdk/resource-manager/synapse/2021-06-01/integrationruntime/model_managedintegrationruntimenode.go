package integrationruntime

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedIntegrationRuntimeNode struct {
	Errors *[]ManagedIntegrationRuntimeError    `json:"errors,omitempty"`
	NodeId *string                              `json:"nodeId,omitempty"`
	Status *ManagedIntegrationRuntimeNodeStatus `json:"status,omitempty"`
}
