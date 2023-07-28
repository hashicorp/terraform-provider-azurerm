package modelversion

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ModelVersion struct {
	Description       *string                 `json:"description,omitempty"`
	Flavors           *map[string]FlavorData  `json:"flavors,omitempty"`
	IsAnonymous       *bool                   `json:"isAnonymous,omitempty"`
	IsArchived        *bool                   `json:"isArchived,omitempty"`
	JobName           *string                 `json:"jobName,omitempty"`
	ModelType         *string                 `json:"modelType,omitempty"`
	ModelUri          *string                 `json:"modelUri,omitempty"`
	Properties        *map[string]string      `json:"properties,omitempty"`
	ProvisioningState *AssetProvisioningState `json:"provisioningState,omitempty"`
	Stage             *string                 `json:"stage,omitempty"`
	Tags              *map[string]string      `json:"tags,omitempty"`
}
