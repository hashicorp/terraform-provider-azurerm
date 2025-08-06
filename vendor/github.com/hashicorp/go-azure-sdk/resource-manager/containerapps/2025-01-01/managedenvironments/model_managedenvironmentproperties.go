package managedenvironments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedEnvironmentProperties struct {
	AppLogsConfiguration        *AppLogsConfiguration                                 `json:"appLogsConfiguration,omitempty"`
	CustomDomainConfiguration   *CustomDomainConfiguration                            `json:"customDomainConfiguration,omitempty"`
	DaprAIConnectionString      *string                                               `json:"daprAIConnectionString,omitempty"`
	DaprAIInstrumentationKey    *string                                               `json:"daprAIInstrumentationKey,omitempty"`
	DaprConfiguration           *DaprConfiguration                                    `json:"daprConfiguration,omitempty"`
	DefaultDomain               *string                                               `json:"defaultDomain,omitempty"`
	DeploymentErrors            *string                                               `json:"deploymentErrors,omitempty"`
	EventStreamEndpoint         *string                                               `json:"eventStreamEndpoint,omitempty"`
	InfrastructureResourceGroup *string                                               `json:"infrastructureResourceGroup,omitempty"`
	KedaConfiguration           *KedaConfiguration                                    `json:"kedaConfiguration,omitempty"`
	PeerAuthentication          *ManagedEnvironmentPropertiesPeerAuthentication       `json:"peerAuthentication,omitempty"`
	PeerTrafficConfiguration    *ManagedEnvironmentPropertiesPeerTrafficConfiguration `json:"peerTrafficConfiguration,omitempty"`
	ProvisioningState           *EnvironmentProvisioningState                         `json:"provisioningState,omitempty"`
	StaticIP                    *string                                               `json:"staticIp,omitempty"`
	VnetConfiguration           *VnetConfiguration                                    `json:"vnetConfiguration,omitempty"`
	WorkloadProfiles            *[]WorkloadProfile                                    `json:"workloadProfiles,omitempty"`
	ZoneRedundant               *bool                                                 `json:"zoneRedundant,omitempty"`
}
