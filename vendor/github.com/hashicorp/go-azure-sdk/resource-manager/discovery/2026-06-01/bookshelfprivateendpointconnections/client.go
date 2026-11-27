package bookshelfprivateendpointconnections

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BookshelfPrivateEndpointConnectionsClient struct {
	Client *resourcemanager.Client
}

func NewBookshelfPrivateEndpointConnectionsClientWithBaseURI(sdkApi sdkEnv.Api) (*BookshelfPrivateEndpointConnectionsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "bookshelfprivateendpointconnections", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating BookshelfPrivateEndpointConnectionsClient: %+v", err)
	}

	return &BookshelfPrivateEndpointConnectionsClient{
		Client: client,
	}, nil
}
