package assets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Event struct {
	EventConfiguration *string                 `json:"eventConfiguration,omitempty"`
	EventNotifier      string                  `json:"eventNotifier"`
	Name               string                  `json:"name"`
	ObservabilityMode  *EventObservabilityMode `json:"observabilityMode,omitempty"`
	Topic              *Topic                  `json:"topic,omitempty"`
}
