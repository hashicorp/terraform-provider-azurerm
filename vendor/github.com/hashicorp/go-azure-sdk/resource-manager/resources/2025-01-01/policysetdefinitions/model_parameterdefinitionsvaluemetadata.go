package policysetdefinitions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ParameterDefinitionsValueMetadata struct {
	AssignPermissions *bool   `json:"assignPermissions,omitempty"`
	Description       *string `json:"description,omitempty"`
	DisplayName       *string `json:"displayName,omitempty"`
	StrongType        *string `json:"strongType,omitempty"`
}
