package privateendpoints

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateEndpointsClient struct {
	Client *resourcemanager.Client
}

func NewPrivateEndpointsClientWithBaseURI(sdkApi sdkEnv.Api) (*PrivateEndpointsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "privateendpoints", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PrivateEndpointsClient: %+v", err)
	}

	return &PrivateEndpointsClient{
		Client: client,
	}, nil
}
