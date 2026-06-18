package privatelinkhubs

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateLinkHubsClient struct {
	Client *resourcemanager.Client
}

func NewPrivateLinkHubsClientWithBaseURI(sdkApi sdkEnv.Api) (*PrivateLinkHubsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "privatelinkhubs", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PrivateLinkHubsClient: %+v", err)
	}

	return &PrivateLinkHubsClient{
		Client: client,
	}, nil
}
