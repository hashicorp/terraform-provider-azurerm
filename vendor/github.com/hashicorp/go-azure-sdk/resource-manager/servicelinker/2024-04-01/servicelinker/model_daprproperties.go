package servicelinker

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DaprProperties struct {
	BindingComponentDirection *DaprBindingComponentDirection `json:"bindingComponentDirection,omitempty"`
	ComponentType             *string                        `json:"componentType,omitempty"`
	Metadata                  *[]DaprMetadata                `json:"metadata,omitempty"`
	RuntimeVersion            *string                        `json:"runtimeVersion,omitempty"`
	Scopes                    *[]string                      `json:"scopes,omitempty"`
	SecretStoreComponent      *string                        `json:"secretStoreComponent,omitempty"`
	Version                   *string                        `json:"version,omitempty"`
}
