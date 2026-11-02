package cognitiveservicesprojects

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProjectProperties struct {
	Description       *string            `json:"description,omitempty"`
	DisplayName       *string            `json:"displayName,omitempty"`
	Endpoints         *map[string]string `json:"endpoints,omitempty"`
	IsDefault         *bool              `json:"isDefault,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
}
