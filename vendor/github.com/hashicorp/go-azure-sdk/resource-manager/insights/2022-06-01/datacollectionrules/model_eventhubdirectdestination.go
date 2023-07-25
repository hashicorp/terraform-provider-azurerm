package datacollectionrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventHubDirectDestination struct {
	EventHubResourceId *string `json:"eventHubResourceId,omitempty"`
	Name               *string `json:"name,omitempty"`
}
