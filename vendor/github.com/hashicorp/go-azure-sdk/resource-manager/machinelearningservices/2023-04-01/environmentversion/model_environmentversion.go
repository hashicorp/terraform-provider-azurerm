package environmentversion

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnvironmentVersion struct {
	AutoRebuild       *AutoRebuildSetting           `json:"autoRebuild,omitempty"`
	Build             *BuildContext                 `json:"build,omitempty"`
	CondaFile         *string                       `json:"condaFile,omitempty"`
	Description       *string                       `json:"description,omitempty"`
	EnvironmentType   *EnvironmentType              `json:"environmentType,omitempty"`
	Image             *string                       `json:"image,omitempty"`
	InferenceConfig   *InferenceContainerProperties `json:"inferenceConfig,omitempty"`
	IsAnonymous       *bool                         `json:"isAnonymous,omitempty"`
	IsArchived        *bool                         `json:"isArchived,omitempty"`
	OsType            *OperatingSystemType          `json:"osType,omitempty"`
	Properties        *map[string]string            `json:"properties,omitempty"`
	ProvisioningState *AssetProvisioningState       `json:"provisioningState,omitempty"`
	Stage             *string                       `json:"stage,omitempty"`
	Tags              *map[string]string            `json:"tags,omitempty"`
}
