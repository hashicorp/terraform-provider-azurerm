package availabilitygrouplisteners

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LoadBalancerConfiguration struct {
	LoadBalancerResourceId     *string           `json:"loadBalancerResourceId,omitempty"`
	PrivateIPAddress           *PrivateIPAddress `json:"privateIpAddress,omitempty"`
	ProbePort                  *int64            `json:"probePort,omitempty"`
	PublicIPAddressResourceId  *string           `json:"publicIpAddressResourceId,omitempty"`
	SqlVirtualMachineInstances *[]string         `json:"sqlVirtualMachineInstances,omitempty"`
}
