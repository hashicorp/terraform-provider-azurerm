package environments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Gen2EnvironmentCreationProperties struct {
	StorageConfiguration   Gen2StorageConfigurationInput     `json:"storageConfiguration"`
	TimeSeriesIdProperties []TimeSeriesIdProperty            `json:"timeSeriesIdProperties"`
	WarmStoreConfiguration *WarmStoreConfigurationProperties `json:"warmStoreConfiguration,omitempty"`
}
