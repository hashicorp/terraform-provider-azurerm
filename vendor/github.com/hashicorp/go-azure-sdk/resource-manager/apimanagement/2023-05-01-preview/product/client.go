package product

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProductClient struct {
	Client *resourcemanager.Client
}

func NewProductClientWithBaseURI(sdkApi sdkEnv.Api) (*ProductClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "product", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ProductClient: %+v", err)
	}

	return &ProductClient{
		Client: client,
	}, nil
}
