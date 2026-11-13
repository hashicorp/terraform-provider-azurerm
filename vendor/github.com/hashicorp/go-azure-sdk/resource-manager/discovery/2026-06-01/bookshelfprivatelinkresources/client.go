package bookshelfprivatelinkresources

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BookshelfPrivateLinkResourcesClient struct {
	Client *resourcemanager.Client
}

func NewBookshelfPrivateLinkResourcesClientWithBaseURI(sdkApi sdkEnv.Api) (*BookshelfPrivateLinkResourcesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "bookshelfprivatelinkresources", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating BookshelfPrivateLinkResourcesClient: %+v", err)
	}

	return &BookshelfPrivateLinkResourcesClient{
		Client: client,
	}, nil
}
