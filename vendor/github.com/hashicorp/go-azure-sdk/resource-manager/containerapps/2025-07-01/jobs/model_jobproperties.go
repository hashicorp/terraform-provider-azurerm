package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobProperties struct {
	Configuration       *JobConfiguration     `json:"configuration,omitempty"`
	EnvironmentId       *string               `json:"environmentId,omitempty"`
	EventStreamEndpoint *string               `json:"eventStreamEndpoint,omitempty"`
	OutboundIPAddresses *[]string             `json:"outboundIpAddresses,omitempty"`
	ProvisioningState   *JobProvisioningState `json:"provisioningState,omitempty"`
	Template            *JobTemplate          `json:"template,omitempty"`
	WorkloadProfileName *string               `json:"workloadProfileName,omitempty"`
}
