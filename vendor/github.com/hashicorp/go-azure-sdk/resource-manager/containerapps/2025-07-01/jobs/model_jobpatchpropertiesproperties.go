package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobPatchPropertiesProperties struct {
	Configuration       *JobConfiguration `json:"configuration,omitempty"`
	EnvironmentId       *string           `json:"environmentId,omitempty"`
	EventStreamEndpoint *string           `json:"eventStreamEndpoint,omitempty"`
	OutboundIPAddresses *[]string         `json:"outboundIpAddresses,omitempty"`
	Template            *JobTemplate      `json:"template,omitempty"`
}
