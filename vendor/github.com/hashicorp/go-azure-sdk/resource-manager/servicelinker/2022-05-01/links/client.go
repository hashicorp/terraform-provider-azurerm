package links

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinksClient struct {
	Client *resourcemanager.Client
}

func NewLinksClientWithBaseURI(sdkApi sdkEnv.Api) (*LinksClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "links", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating LinksClient: %+v", err)
	}

	return &LinksClient{
		Client: client,
	}, nil
}
