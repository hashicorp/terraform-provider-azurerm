package componentversion

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ComponentVersion struct {
	ComponentSpec     *interface{}            `json:"componentSpec,omitempty"`
	Description       *string                 `json:"description,omitempty"`
	IsAnonymous       *bool                   `json:"isAnonymous,omitempty"`
	IsArchived        *bool                   `json:"isArchived,omitempty"`
	Properties        *map[string]string      `json:"properties,omitempty"`
	ProvisioningState *AssetProvisioningState `json:"provisioningState,omitempty"`
	Tags              *map[string]string      `json:"tags,omitempty"`
}
