package vpnserverconfigurations

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VpnServerConfigurationsClient struct {
	Client *resourcemanager.Client
}

func NewVpnServerConfigurationsClientWithBaseURI(api environments.Api) (*VpnServerConfigurationsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "vpnserverconfigurations", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating VpnServerConfigurationsClient: %+v", err)
	}

	return &VpnServerConfigurationsClient{
		Client: client,
	}, nil
}
