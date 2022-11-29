package sqlvirtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerConfigurationsManagementSettings struct {
	AdditionalFeaturesServerConfigurations *AdditionalFeaturesServerConfigurations `json:"additionalFeaturesServerConfigurations,omitempty"`
	SqlConnectivityUpdateSettings          *SqlConnectivityUpdateSettings          `json:"sqlConnectivityUpdateSettings,omitempty"`
	SqlInstanceSettings                    *SQLInstanceSettings                    `json:"sqlInstanceSettings,omitempty"`
	SqlStorageUpdateSettings               *SqlStorageUpdateSettings               `json:"sqlStorageUpdateSettings,omitempty"`
	SqlWorkloadTypeUpdateSettings          *SqlWorkloadTypeUpdateSettings          `json:"sqlWorkloadTypeUpdateSettings,omitempty"`
}
