package hcpopenshiftclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkProfile struct {
	HostPrefix  *int64       `json:"hostPrefix,omitempty"`
	MachineCidr *string      `json:"machineCidr,omitempty"`
	NetworkType *NetworkType `json:"networkType,omitempty"`
	PodCidr     *string      `json:"podCidr,omitempty"`
	ServiceCidr *string      `json:"serviceCidr,omitempty"`
}
