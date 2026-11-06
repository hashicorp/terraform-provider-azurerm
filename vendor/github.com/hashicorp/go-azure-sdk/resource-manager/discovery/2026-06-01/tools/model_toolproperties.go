package tools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ToolProperties struct {
	DefinitionContent    map[string]interface{} `json:"definitionContent"`
	EnvironmentVariables *map[string]string     `json:"environmentVariables,omitempty"`
	ProvisioningState    *ProvisioningState     `json:"provisioningState,omitempty"`
	Version              string                 `json:"version"`
}
