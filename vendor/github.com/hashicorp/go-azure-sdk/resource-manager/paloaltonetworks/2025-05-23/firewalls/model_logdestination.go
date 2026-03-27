package firewalls

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LogDestination struct {
	EventHubConfigurations *EventHub       `json:"eventHubConfigurations,omitempty"`
	MonitorConfigurations  *MonitorLog     `json:"monitorConfigurations,omitempty"`
	StorageConfigurations  *StorageAccount `json:"storageConfigurations,omitempty"`
}
