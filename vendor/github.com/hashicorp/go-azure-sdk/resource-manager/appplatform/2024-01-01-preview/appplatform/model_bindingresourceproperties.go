package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BindingResourceProperties struct {
	BindingParameters   *map[string]string `json:"bindingParameters,omitempty"`
	CreatedAt           *string            `json:"createdAt,omitempty"`
	GeneratedProperties *string            `json:"generatedProperties,omitempty"`
	Key                 *string            `json:"key,omitempty"`
	ResourceId          *string            `json:"resourceId,omitempty"`
	ResourceName        *string            `json:"resourceName,omitempty"`
	ResourceType        *string            `json:"resourceType,omitempty"`
	UpdatedAt           *string            `json:"updatedAt,omitempty"`
}
