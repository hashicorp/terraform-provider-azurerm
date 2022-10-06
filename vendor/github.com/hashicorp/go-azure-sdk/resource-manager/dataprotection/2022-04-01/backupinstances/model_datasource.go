package backupinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Datasource struct {
	DatasourceType   *string `json:"datasourceType,omitempty"`
	ObjectType       *string `json:"objectType,omitempty"`
	ResourceID       string  `json:"resourceID"`
	ResourceLocation *string `json:"resourceLocation,omitempty"`
	ResourceName     *string `json:"resourceName,omitempty"`
	ResourceType     *string `json:"resourceType,omitempty"`
	ResourceUri      *string `json:"resourceUri,omitempty"`
}
