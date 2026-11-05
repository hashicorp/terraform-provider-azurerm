// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	network_2025_07_01 "github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecurityperimeteraccessrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecurityperimeterassociations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecurityperimeterprofiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecurityperimeters"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	*network_2025_07_01.Client

	NetworkSecurityPerimeterAccessRulesClient  *networksecurityperimeteraccessrules.NetworkSecurityPerimeterAccessRulesClient
	NetworkSecurityPerimeterAssociationsClient *networksecurityperimeterassociations.NetworkSecurityPerimeterAssociationsClient
	NetworkSecurityPerimeterProfilesClient     *networksecurityperimeterprofiles.NetworkSecurityPerimeterProfilesClient
	NetworkSecurityPerimetersClient            *networksecurityperimeters.NetworkSecurityPerimetersClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	NetworkSecurityPerimeterAssociationsClient, err := networksecurityperimeterassociations.NewNetworkSecurityPerimeterAssociationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Network Security Perimeter Resource Association Client: %+v", err)
	}
	o.Configure(NetworkSecurityPerimeterAssociationsClient.Client, o.Authorizers.ResourceManager)

	NetworkSecurityPerimeterAccessRulesClient, err := networksecurityperimeteraccessrules.NewNetworkSecurityPerimeterAccessRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Network Security Perimeter Access Rules Client: %+v", err)
	}
	o.Configure(NetworkSecurityPerimeterAccessRulesClient.Client, o.Authorizers.ResourceManager)

	NetworkSecurityPerimeterProfilesClient, err := networksecurityperimeterprofiles.NewNetworkSecurityPerimeterProfilesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Network Security Perimeter Profiles Client: %+v", err)
	}
	o.Configure(NetworkSecurityPerimeterProfilesClient.Client, o.Authorizers.ResourceManager)

	NetworkSecurityPerimetersClient, err := networksecurityperimeters.NewNetworkSecurityPerimetersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Network Security Perimeters Client: %+v", err)
	}
	o.Configure(NetworkSecurityPerimetersClient.Client, o.Authorizers.ResourceManager)

	client, err := network_2025_07_01.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building clients for Network Manager: %+v", err)
	}

	return &Client{
		NetworkSecurityPerimeterAccessRulesClient:  NetworkSecurityPerimeterAccessRulesClient,
		NetworkSecurityPerimeterAssociationsClient: NetworkSecurityPerimeterAssociationsClient,
		NetworkSecurityPerimeterProfilesClient:     NetworkSecurityPerimeterProfilesClient,
		NetworkSecurityPerimetersClient:            NetworkSecurityPerimetersClient,
		Client:                                     client,
	}, nil
}
