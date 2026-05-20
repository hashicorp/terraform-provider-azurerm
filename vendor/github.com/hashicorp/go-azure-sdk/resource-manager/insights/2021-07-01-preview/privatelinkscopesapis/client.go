package privatelinkscopesapis

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateLinkScopesAPIsClient struct {
	Client *resourcemanager.Client
}

func NewPrivateLinkScopesAPIsClientWithBaseURI(sdkApi sdkEnv.Api) (*PrivateLinkScopesAPIsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "privatelinkscopesapis", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PrivateLinkScopesAPIsClient: %+v", err)
	}

	return &PrivateLinkScopesAPIsClient{
		Client: client,
	}, nil
}
