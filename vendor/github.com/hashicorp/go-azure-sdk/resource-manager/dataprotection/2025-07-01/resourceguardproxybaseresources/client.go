package resourceguardproxybaseresources

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceGuardProxyBaseResourcesClient struct {
	Client *resourcemanager.Client
}

func NewResourceGuardProxyBaseResourcesClientWithBaseURI(sdkApi sdkEnv.Api) (*ResourceGuardProxyBaseResourcesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "resourceguardproxybaseresources", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ResourceGuardProxyBaseResourcesClient: %+v", err)
	}

	return &ResourceGuardProxyBaseResourcesClient{
		Client: client,
	}, nil
}
