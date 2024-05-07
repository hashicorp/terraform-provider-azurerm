package privatelinkservice

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateLinkServiceClient struct {
	Client *resourcemanager.Client
}

func NewPrivateLinkServiceClientWithBaseURI(sdkApi sdkEnv.Api) (*PrivateLinkServiceClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "privatelinkservice", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PrivateLinkServiceClient: %+v", err)
	}

	return &PrivateLinkServiceClient{
		Client: client,
	}, nil
}
