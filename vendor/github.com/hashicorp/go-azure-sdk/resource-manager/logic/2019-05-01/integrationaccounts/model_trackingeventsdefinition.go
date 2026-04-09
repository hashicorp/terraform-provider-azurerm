package integrationaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TrackingEventsDefinition struct {
	Events             []TrackingEvent              `json:"events"`
	SourceType         string                       `json:"sourceType"`
	TrackEventsOptions *TrackEventsOperationOptions `json:"trackEventsOptions,omitempty"`
}
