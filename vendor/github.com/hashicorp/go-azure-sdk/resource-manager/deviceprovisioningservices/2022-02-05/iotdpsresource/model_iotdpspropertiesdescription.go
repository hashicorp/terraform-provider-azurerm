package iotdpsresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IotDpsPropertiesDescription struct {
	AllocationPolicy           *AllocationPolicy                                                `json:"allocationPolicy,omitempty"`
	AuthorizationPolicies      *[]SharedAccessSignatureAuthorizationRuleAccessRightsDescription `json:"authorizationPolicies,omitempty"`
	DeviceProvisioningHostName *string                                                          `json:"deviceProvisioningHostName,omitempty"`
	EnableDataResidency        *bool                                                            `json:"enableDataResidency,omitempty"`
	IPFilterRules              *[]IPFilterRule                                                  `json:"ipFilterRules,omitempty"`
	IdScope                    *string                                                          `json:"idScope,omitempty"`
	IotHubs                    *[]IotHubDefinitionDescription                                   `json:"iotHubs,omitempty"`
	PrivateEndpointConnections *[]PrivateEndpointConnection                                     `json:"privateEndpointConnections,omitempty"`
	ProvisioningState          *string                                                          `json:"provisioningState,omitempty"`
	PublicNetworkAccess        *PublicNetworkAccess                                             `json:"publicNetworkAccess,omitempty"`
	ServiceOperationsHostName  *string                                                          `json:"serviceOperationsHostName,omitempty"`
	State                      *State                                                           `json:"state,omitempty"`
}
