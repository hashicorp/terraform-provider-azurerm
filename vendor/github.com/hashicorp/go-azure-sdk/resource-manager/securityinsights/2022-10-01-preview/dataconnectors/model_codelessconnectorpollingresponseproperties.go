package dataconnectors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CodelessConnectorPollingResponseProperties struct {
	EventsJsonPaths       []string `json:"eventsJsonPaths"`
	IsGzipCompressed      *bool    `json:"isGzipCompressed,omitempty"`
	SuccessStatusJsonPath *string  `json:"successStatusJsonPath,omitempty"`
	SuccessStatusValue    *string  `json:"successStatusValue,omitempty"`
}
