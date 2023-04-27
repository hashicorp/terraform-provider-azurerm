package assignment

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssignmentProperties struct {
	BlueprintId       *string                       `json:"blueprintId,omitempty"`
	Description       *string                       `json:"description,omitempty"`
	DisplayName       *string                       `json:"displayName,omitempty"`
	Locks             *AssignmentLockSettings       `json:"locks,omitempty"`
	Parameters        map[string]ParameterValue     `json:"parameters"`
	ProvisioningState *AssignmentProvisioningState  `json:"provisioningState,omitempty"`
	ResourceGroups    map[string]ResourceGroupValue `json:"resourceGroups"`
	Scope             *string                       `json:"scope,omitempty"`
	Status            *AssignmentStatus             `json:"status,omitempty"`
}
