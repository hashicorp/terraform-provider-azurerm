package policysetdefinitions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PolicyDefinitionGroup struct {
	AdditionalMetadataId *string `json:"additionalMetadataId,omitempty"`
	Category             *string `json:"category,omitempty"`
	Description          *string `json:"description,omitempty"`
	DisplayName          *string `json:"displayName,omitempty"`
	Name                 string  `json:"name"`
}
