package vpnlinkconnections

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VpnLinkConnectionsClient struct {
	Client *resourcemanager.Client
}

func NewVpnLinkConnectionsClientWithBaseURI(sdkApi sdkEnv.Api) (*VpnLinkConnectionsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "vpnlinkconnections", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating VpnLinkConnectionsClient: %+v", err)
	}

	return &VpnLinkConnectionsClient{
		Client: client,
	}, nil
}
