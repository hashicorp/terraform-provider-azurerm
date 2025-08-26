package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AksNetworkingConfiguration struct {
	DnsServiceIP     *string `json:"dnsServiceIP,omitempty"`
	DockerBridgeCidr *string `json:"dockerBridgeCidr,omitempty"`
	ServiceCidr      *string `json:"serviceCidr,omitempty"`
	SubnetId         *string `json:"subnetId,omitempty"`
}
