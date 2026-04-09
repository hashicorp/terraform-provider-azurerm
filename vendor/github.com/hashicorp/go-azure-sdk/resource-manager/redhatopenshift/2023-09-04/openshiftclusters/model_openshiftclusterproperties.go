package openshiftclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OpenShiftClusterProperties struct {
	ApiserverProfile        *APIServerProfile        `json:"apiserverProfile,omitempty"`
	ClusterProfile          *ClusterProfile          `json:"clusterProfile,omitempty"`
	ConsoleProfile          *ConsoleProfile          `json:"consoleProfile,omitempty"`
	IngressProfiles         *[]IngressProfile        `json:"ingressProfiles,omitempty"`
	MasterProfile           *MasterProfile           `json:"masterProfile,omitempty"`
	NetworkProfile          *NetworkProfile          `json:"networkProfile,omitempty"`
	ProvisioningState       *ProvisioningState       `json:"provisioningState,omitempty"`
	ServicePrincipalProfile *ServicePrincipalProfile `json:"servicePrincipalProfile,omitempty"`
	WorkerProfiles          *[]WorkerProfile         `json:"workerProfiles,omitempty"`
	WorkerProfilesStatus    *[]WorkerProfile         `json:"workerProfilesStatus,omitempty"`
}
