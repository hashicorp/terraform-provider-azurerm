package projectpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourcePolicy struct {
	Action       *PolicyAction          `json:"action,omitempty"`
	Filter       *string                `json:"filter,omitempty"`
	ResourceType *DevCenterResourceType `json:"resourceType,omitempty"`
	Resources    *string                `json:"resources,omitempty"`
}
