package publishedblueprint

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PublishedBlueprintProperties struct {
	BlueprintName  *string                             `json:"blueprintName,omitempty"`
	ChangeNotes    *string                             `json:"changeNotes,omitempty"`
	Description    *string                             `json:"description,omitempty"`
	DisplayName    *string                             `json:"displayName,omitempty"`
	Parameters     *map[string]ParameterDefinition     `json:"parameters,omitempty"`
	ResourceGroups *map[string]ResourceGroupDefinition `json:"resourceGroups,omitempty"`
	Status         *BlueprintResourceStatusBase        `json:"status,omitempty"`
	TargetScope    *BlueprintTargetScope               `json:"targetScope,omitempty"`
}
