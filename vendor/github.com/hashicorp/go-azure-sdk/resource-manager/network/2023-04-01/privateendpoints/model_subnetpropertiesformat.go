package privateendpoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubnetPropertiesFormat struct {
	AddressPrefix                      *string                                          `json:"addressPrefix,omitempty"`
	AddressPrefixes                    *[]string                                        `json:"addressPrefixes,omitempty"`
	ApplicationGatewayIPConfigurations *[]ApplicationGatewayIPConfiguration             `json:"applicationGatewayIPConfigurations,omitempty"`
	Delegations                        *[]Delegation                                    `json:"delegations,omitempty"`
	IPAllocations                      *[]SubResource                                   `json:"ipAllocations,omitempty"`
	IPConfigurationProfiles            *[]IPConfigurationProfile                        `json:"ipConfigurationProfiles,omitempty"`
	IPConfigurations                   *[]IPConfiguration                               `json:"ipConfigurations,omitempty"`
	NatGateway                         *SubResource                                     `json:"natGateway,omitempty"`
	NetworkSecurityGroup               *NetworkSecurityGroup                            `json:"networkSecurityGroup,omitempty"`
	PrivateEndpointNetworkPolicies     *VirtualNetworkPrivateEndpointNetworkPolicies    `json:"privateEndpointNetworkPolicies,omitempty"`
	PrivateEndpoints                   *[]PrivateEndpoint                               `json:"privateEndpoints,omitempty"`
	PrivateLinkServiceNetworkPolicies  *VirtualNetworkPrivateLinkServiceNetworkPolicies `json:"privateLinkServiceNetworkPolicies,omitempty"`
	ProvisioningState                  *ProvisioningState                               `json:"provisioningState,omitempty"`
	Purpose                            *string                                          `json:"purpose,omitempty"`
	ResourceNavigationLinks            *[]ResourceNavigationLink                        `json:"resourceNavigationLinks,omitempty"`
	RouteTable                         *RouteTable                                      `json:"routeTable,omitempty"`
	ServiceAssociationLinks            *[]ServiceAssociationLink                        `json:"serviceAssociationLinks,omitempty"`
	ServiceEndpointPolicies            *[]ServiceEndpointPolicy                         `json:"serviceEndpointPolicies,omitempty"`
	ServiceEndpoints                   *[]ServiceEndpointPropertiesFormat               `json:"serviceEndpoints,omitempty"`
}
