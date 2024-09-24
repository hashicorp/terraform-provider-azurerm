package sku

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SkuClient struct {
	Client *resourcemanager.Client
}

func NewSkuClientWithBaseURI(sdkApi sdkEnv.Api) (*SkuClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "sku", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SkuClient: %+v", err)
	}

	return &SkuClient{
		Client: client,
	}, nil
}
