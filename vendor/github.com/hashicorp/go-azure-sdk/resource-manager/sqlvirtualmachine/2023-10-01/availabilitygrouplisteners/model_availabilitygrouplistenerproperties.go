package availabilitygrouplisteners

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AvailabilityGroupListenerProperties struct {
	AvailabilityGroupConfiguration           *AgConfiguration              `json:"availabilityGroupConfiguration,omitempty"`
	AvailabilityGroupName                    *string                       `json:"availabilityGroupName,omitempty"`
	CreateDefaultAvailabilityGroupIfNotExist *bool                         `json:"createDefaultAvailabilityGroupIfNotExist,omitempty"`
	LoadBalancerConfigurations               *[]LoadBalancerConfiguration  `json:"loadBalancerConfigurations,omitempty"`
	MultiSubnetIPConfigurations              *[]MultiSubnetIPConfiguration `json:"multiSubnetIpConfigurations,omitempty"`
	Port                                     *int64                        `json:"port,omitempty"`
	ProvisioningState                        *string                       `json:"provisioningState,omitempty"`
}
