package streamingjobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventGridStreamInputDataSourceProperties struct {
	EventTypes      *[]string                        `json:"eventTypes,omitempty"`
	Schema          *EventGridEventSchemaType        `json:"schema,omitempty"`
	StorageAccounts *[]StorageAccount                `json:"storageAccounts,omitempty"`
	Subscriber      *EventHubV2StreamInputDataSource `json:"subscriber,omitempty"`
}
