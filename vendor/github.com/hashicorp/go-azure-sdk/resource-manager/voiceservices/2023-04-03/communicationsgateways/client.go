package communicationsgateways

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CommunicationsGatewaysClient struct {
	Client *resourcemanager.Client
}

func NewCommunicationsGatewaysClientWithBaseURI(sdkApi sdkEnv.Api) (*CommunicationsGatewaysClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "communicationsgateways", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating CommunicationsGatewaysClient: %+v", err)
	}

	return &CommunicationsGatewaysClient{
		Client: client,
	}, nil
}
