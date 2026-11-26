package virtualnetworkappliancecapabilities

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkApplianceCapabilitiesClient struct {
	Client *resourcemanager.Client
}

func NewVirtualNetworkApplianceCapabilitiesClientWithBaseURI(sdkApi sdkEnv.Api) (*VirtualNetworkApplianceCapabilitiesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "virtualnetworkappliancecapabilities", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating VirtualNetworkApplianceCapabilitiesClient: %+v", err)
	}

	return &VirtualNetworkApplianceCapabilitiesClient{
		Client: client,
	}, nil
}
