package productgroup

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProductGroupClient struct {
	Client *resourcemanager.Client
}

func NewProductGroupClientWithBaseURI(sdkApi sdkEnv.Api) (*ProductGroupClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "productgroup", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ProductGroupClient: %+v", err)
	}

	return &ProductGroupClient{
		Client: client,
	}, nil
}
