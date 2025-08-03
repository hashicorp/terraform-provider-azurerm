package datacollectionrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataImportSources struct {
	EventHub *EventHubDataSource `json:"eventHub,omitempty"`
}
