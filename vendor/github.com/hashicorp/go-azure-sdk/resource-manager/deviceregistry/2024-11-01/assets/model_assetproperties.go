package assets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssetProperties struct {
	AssetEndpointProfileRef      string                  `json:"assetEndpointProfileRef"`
	Attributes                   *map[string]interface{} `json:"attributes,omitempty"`
	Datasets                     *[]Dataset              `json:"datasets,omitempty"`
	DefaultDatasetsConfiguration *string                 `json:"defaultDatasetsConfiguration,omitempty"`
	DefaultEventsConfiguration   *string                 `json:"defaultEventsConfiguration,omitempty"`
	DefaultTopic                 *Topic                  `json:"defaultTopic,omitempty"`
	Description                  *string                 `json:"description,omitempty"`
	DiscoveredAssetRefs          *[]string               `json:"discoveredAssetRefs,omitempty"`
	DisplayName                  *string                 `json:"displayName,omitempty"`
	DocumentationUri             *string                 `json:"documentationUri,omitempty"`
	Enabled                      *bool                   `json:"enabled,omitempty"`
	Events                       *[]Event                `json:"events,omitempty"`
	ExternalAssetId              *string                 `json:"externalAssetId,omitempty"`
	HardwareRevision             *string                 `json:"hardwareRevision,omitempty"`
	Manufacturer                 *string                 `json:"manufacturer,omitempty"`
	ManufacturerUri              *string                 `json:"manufacturerUri,omitempty"`
	Model                        *string                 `json:"model,omitempty"`
	ProductCode                  *string                 `json:"productCode,omitempty"`
	ProvisioningState            *ProvisioningState      `json:"provisioningState,omitempty"`
	SerialNumber                 *string                 `json:"serialNumber,omitempty"`
	SoftwareRevision             *string                 `json:"softwareRevision,omitempty"`
	Status                       *AssetStatus            `json:"status,omitempty"`
	Uuid                         *string                 `json:"uuid,omitempty"`
	Version                      *int64                  `json:"version,omitempty"`
}
