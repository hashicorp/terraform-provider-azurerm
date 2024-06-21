package connectivityconfigurations

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectivityConfigurationsClient struct {
	Client *resourcemanager.Client
}

func NewConnectivityConfigurationsClientWithBaseURI(sdkApi sdkEnv.Api) (*ConnectivityConfigurationsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "connectivityconfigurations", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ConnectivityConfigurationsClient: %+v", err)
	}

	return &ConnectivityConfigurationsClient{
		Client: client,
	}, nil
}
