package publicipprefixes

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PublicIPPrefixesClient struct {
	Client *resourcemanager.Client
}

func NewPublicIPPrefixesClientWithBaseURI(sdkApi sdkEnv.Api) (*PublicIPPrefixesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "publicipprefixes", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PublicIPPrefixesClient: %+v", err)
	}

	return &PublicIPPrefixesClient{
		Client: client,
	}, nil
}
