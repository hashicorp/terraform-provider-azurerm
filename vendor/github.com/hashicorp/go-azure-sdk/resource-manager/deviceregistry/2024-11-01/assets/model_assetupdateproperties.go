package assets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssetUpdateProperties struct {
	Attributes                   *map[string]interface{} `json:"attributes,omitempty"`
	Datasets                     *[]Dataset              `json:"datasets,omitempty"`
	DefaultDatasetsConfiguration *string                 `json:"defaultDatasetsConfiguration,omitempty"`
	DefaultEventsConfiguration   *string                 `json:"defaultEventsConfiguration,omitempty"`
	DefaultTopic                 *TopicUpdate            `json:"defaultTopic,omitempty"`
	Description                  *string                 `json:"description,omitempty"`
	DisplayName                  *string                 `json:"displayName,omitempty"`
	DocumentationUri             *string                 `json:"documentationUri,omitempty"`
	Enabled                      *bool                   `json:"enabled,omitempty"`
	Events                       *[]Event                `json:"events,omitempty"`
	HardwareRevision             *string                 `json:"hardwareRevision,omitempty"`
	Manufacturer                 *string                 `json:"manufacturer,omitempty"`
	ManufacturerUri              *string                 `json:"manufacturerUri,omitempty"`
	Model                        *string                 `json:"model,omitempty"`
	ProductCode                  *string                 `json:"productCode,omitempty"`
	SerialNumber                 *string                 `json:"serialNumber,omitempty"`
	SoftwareRevision             *string                 `json:"softwareRevision,omitempty"`
}
