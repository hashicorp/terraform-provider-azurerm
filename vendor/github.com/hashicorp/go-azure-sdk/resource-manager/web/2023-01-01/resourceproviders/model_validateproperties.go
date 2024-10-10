package resourceproviders

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ValidateProperties struct {
	AppServiceEnvironment     *AppServiceEnvironment `json:"appServiceEnvironment,omitempty"`
	Capacity                  *int64                 `json:"capacity,omitempty"`
	ContainerImagePlatform    *string                `json:"containerImagePlatform,omitempty"`
	ContainerImageRepository  *string                `json:"containerImageRepository,omitempty"`
	ContainerImageTag         *string                `json:"containerImageTag,omitempty"`
	ContainerRegistryBaseURL  *string                `json:"containerRegistryBaseUrl,omitempty"`
	ContainerRegistryPassword *string                `json:"containerRegistryPassword,omitempty"`
	ContainerRegistryUsername *string                `json:"containerRegistryUsername,omitempty"`
	HostingEnvironment        *string                `json:"hostingEnvironment,omitempty"`
	IsSpot                    *bool                  `json:"isSpot,omitempty"`
	IsXenon                   *bool                  `json:"isXenon,omitempty"`
	NeedLinuxWorkers          *bool                  `json:"needLinuxWorkers,omitempty"`
	ServerFarmId              *string                `json:"serverFarmId,omitempty"`
	SkuName                   *string                `json:"skuName,omitempty"`
}
