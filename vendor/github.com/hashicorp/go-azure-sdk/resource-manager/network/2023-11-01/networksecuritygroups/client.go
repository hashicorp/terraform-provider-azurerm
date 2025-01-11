package networksecuritygroups

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkSecurityGroupsClient struct {
	Client *resourcemanager.Client
}

func NewNetworkSecurityGroupsClientWithBaseURI(sdkApi sdkEnv.Api) (*NetworkSecurityGroupsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "networksecuritygroups", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NetworkSecurityGroupsClient: %+v", err)
	}

	return &NetworkSecurityGroupsClient{
		Client: client,
	}, nil
}
