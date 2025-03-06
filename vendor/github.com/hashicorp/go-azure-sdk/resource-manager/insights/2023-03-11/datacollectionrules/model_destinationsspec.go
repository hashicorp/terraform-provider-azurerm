package datacollectionrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DestinationsSpec struct {
	AzureDataExplorer   *[]AdxDestination               `json:"azureDataExplorer,omitempty"`
	AzureMonitorMetrics *AzureMonitorMetricsDestination `json:"azureMonitorMetrics,omitempty"`
	EventHubs           *[]EventHubDestination          `json:"eventHubs,omitempty"`
	EventHubsDirect     *[]EventHubDirectDestination    `json:"eventHubsDirect,omitempty"`
	LogAnalytics        *[]LogAnalyticsDestination      `json:"logAnalytics,omitempty"`
	MicrosoftFabric     *[]MicrosoftFabricDestination   `json:"microsoftFabric,omitempty"`
	MonitoringAccounts  *[]MonitoringAccountDestination `json:"monitoringAccounts,omitempty"`
	StorageAccounts     *[]StorageBlobDestination       `json:"storageAccounts,omitempty"`
	StorageBlobsDirect  *[]StorageBlobDestination       `json:"storageBlobsDirect,omitempty"`
	StorageTablesDirect *[]StorageTableDestination      `json:"storageTablesDirect,omitempty"`
}
