package connectedregistries

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectedRegistriesClient struct {
	Client *resourcemanager.Client
}

func NewConnectedRegistriesClientWithBaseURI(sdkApi sdkEnv.Api) (*ConnectedRegistriesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "connectedregistries", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ConnectedRegistriesClient: %+v", err)
	}

	return &ConnectedRegistriesClient{
		Client: client,
	}, nil
}
