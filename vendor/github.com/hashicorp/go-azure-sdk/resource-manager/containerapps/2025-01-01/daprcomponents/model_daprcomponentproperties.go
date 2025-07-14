package daprcomponents

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DaprComponentProperties struct {
	ComponentType        *string         `json:"componentType,omitempty"`
	IgnoreErrors         *bool           `json:"ignoreErrors,omitempty"`
	InitTimeout          *string         `json:"initTimeout,omitempty"`
	Metadata             *[]DaprMetadata `json:"metadata,omitempty"`
	Scopes               *[]string       `json:"scopes,omitempty"`
	SecretStoreComponent *string         `json:"secretStoreComponent,omitempty"`
	Secrets              *[]Secret       `json:"secrets,omitempty"`
	Version              *string         `json:"version,omitempty"`
}
