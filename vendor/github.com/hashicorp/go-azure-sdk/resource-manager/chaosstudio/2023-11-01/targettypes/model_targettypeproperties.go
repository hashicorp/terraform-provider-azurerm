package targettypes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TargetTypeProperties struct {
	Description      *string   `json:"description,omitempty"`
	DisplayName      *string   `json:"displayName,omitempty"`
	PropertiesSchema *string   `json:"propertiesSchema,omitempty"`
	ResourceTypes    *[]string `json:"resourceTypes,omitempty"`
}
