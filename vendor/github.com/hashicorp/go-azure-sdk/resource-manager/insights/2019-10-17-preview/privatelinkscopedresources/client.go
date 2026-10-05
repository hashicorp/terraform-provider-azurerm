package privatelinkscopedresources

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateLinkScopedResourcesClient struct {
	Client *resourcemanager.Client
}

func NewPrivateLinkScopedResourcesClientWithBaseURI(sdkApi sdkEnv.Api) (*PrivateLinkScopedResourcesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "privatelinkscopedresources", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PrivateLinkScopedResourcesClient: %+v", err)
	}

	return &PrivateLinkScopedResourcesClient{
		Client: client,
	}, nil
}
