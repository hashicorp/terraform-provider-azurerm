package prefixlistresources

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrefixListResourcesClient struct {
	Client *resourcemanager.Client
}

func NewPrefixListResourcesClientWithBaseURI(sdkApi sdkEnv.Api) (*PrefixListResourcesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "prefixlistresources", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PrefixListResourcesClient: %+v", err)
	}

	return &PrefixListResourcesClient{
		Client: client,
	}, nil
}
