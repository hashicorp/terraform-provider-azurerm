package servicelinker

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DaprMetadata struct {
	Description *string               `json:"description,omitempty"`
	Name        *string               `json:"name,omitempty"`
	Required    *DaprMetadataRequired `json:"required,omitempty"`
	SecretRef   *string               `json:"secretRef,omitempty"`
	Value       *string               `json:"value,omitempty"`
}
