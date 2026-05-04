package v2024_07_10

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-07-10/extensions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-07-10/licenseprofiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-07-10/licenses"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-07-10/machineextensions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-07-10/machineextensionsupgrade"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-07-10/machinenetworkprofile"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-07-10/machines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-07-10/networksecurityperimeterconfiguration"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-07-10/privateendpointconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-07-10/privatelinkresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-07-10/privatelinkscopes"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	Extensions                            *extensions.ExtensionsClient
	LicenseProfiles                       *licenseprofiles.LicenseProfilesClient
	Licenses                              *licenses.LicensesClient
	MachineExtensions                     *machineextensions.MachineExtensionsClient
	MachineExtensionsUpgrade              *machineextensionsupgrade.MachineExtensionsUpgradeClient
	MachineNetworkProfile                 *machinenetworkprofile.MachineNetworkProfileClient
	Machines                              *machines.MachinesClient
	NetworkSecurityPerimeterConfiguration *networksecurityperimeterconfiguration.NetworkSecurityPerimeterConfigurationClient
	PrivateEndpointConnections            *privateendpointconnections.PrivateEndpointConnectionsClient
	PrivateLinkResources                  *privatelinkresources.PrivateLinkResourcesClient
	PrivateLinkScopes                     *privatelinkscopes.PrivateLinkScopesClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	extensionsClient, err := extensions.NewExtensionsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Extensions client: %+v", err)
	}
	configureFunc(extensionsClient.Client)

	licenseProfilesClient, err := licenseprofiles.NewLicenseProfilesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building LicenseProfiles client: %+v", err)
	}
	configureFunc(licenseProfilesClient.Client)

	licensesClient, err := licenses.NewLicensesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Licenses client: %+v", err)
	}
	configureFunc(licensesClient.Client)

	machineExtensionsClient, err := machineextensions.NewMachineExtensionsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building MachineExtensions client: %+v", err)
	}
	configureFunc(machineExtensionsClient.Client)

	machineExtensionsUpgradeClient, err := machineextensionsupgrade.NewMachineExtensionsUpgradeClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building MachineExtensionsUpgrade client: %+v", err)
	}
	configureFunc(machineExtensionsUpgradeClient.Client)

	machineNetworkProfileClient, err := machinenetworkprofile.NewMachineNetworkProfileClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building MachineNetworkProfile client: %+v", err)
	}
	configureFunc(machineNetworkProfileClient.Client)

	machinesClient, err := machines.NewMachinesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Machines client: %+v", err)
	}
	configureFunc(machinesClient.Client)

	networkSecurityPerimeterConfigurationClient, err := networksecurityperimeterconfiguration.NewNetworkSecurityPerimeterConfigurationClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building NetworkSecurityPerimeterConfiguration client: %+v", err)
	}
	configureFunc(networkSecurityPerimeterConfigurationClient.Client)

	privateEndpointConnectionsClient, err := privateendpointconnections.NewPrivateEndpointConnectionsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PrivateEndpointConnections client: %+v", err)
	}
	configureFunc(privateEndpointConnectionsClient.Client)

	privateLinkResourcesClient, err := privatelinkresources.NewPrivateLinkResourcesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PrivateLinkResources client: %+v", err)
	}
	configureFunc(privateLinkResourcesClient.Client)

	privateLinkScopesClient, err := privatelinkscopes.NewPrivateLinkScopesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PrivateLinkScopes client: %+v", err)
	}
	configureFunc(privateLinkScopesClient.Client)

	return &Client{
		Extensions:                            extensionsClient,
		LicenseProfiles:                       licenseProfilesClient,
		Licenses:                              licensesClient,
		MachineExtensions:                     machineExtensionsClient,
		MachineExtensionsUpgrade:              machineExtensionsUpgradeClient,
		MachineNetworkProfile:                 machineNetworkProfileClient,
		Machines:                              machinesClient,
		NetworkSecurityPerimeterConfiguration: networkSecurityPerimeterConfigurationClient,
		PrivateEndpointConnections:            privateEndpointConnectionsClient,
		PrivateLinkResources:                  privateLinkResourcesClient,
		PrivateLinkScopes:                     privateLinkScopesClient,
	}, nil
}
