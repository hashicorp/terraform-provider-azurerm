package managedenvironments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VnetConfiguration struct {
	DockerBridgeCidr       *string `json:"dockerBridgeCidr,omitempty"`
	InfrastructureSubnetId *string `json:"infrastructureSubnetId,omitempty"`
	Internal               *bool   `json:"internal,omitempty"`
	PlatformReservedCidr   *string `json:"platformReservedCidr,omitempty"`
	PlatformReservedDnsIP  *string `json:"platformReservedDnsIP,omitempty"`
}
