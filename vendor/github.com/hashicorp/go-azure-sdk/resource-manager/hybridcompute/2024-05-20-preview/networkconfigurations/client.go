package networkconfigurations

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkConfigurationsClient struct {
	Client *resourcemanager.Client
}

func NewNetworkConfigurationsClientWithBaseURI(sdkApi sdkEnv.Api) (*NetworkConfigurationsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "networkconfigurations", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NetworkConfigurationsClient: %+v", err)
	}

	return &NetworkConfigurationsClient{
		Client: client,
	}, nil
}
