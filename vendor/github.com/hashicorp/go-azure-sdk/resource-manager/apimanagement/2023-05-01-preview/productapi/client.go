package productapi

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProductApiClient struct {
	Client *resourcemanager.Client
}

func NewProductApiClientWithBaseURI(sdkApi sdkEnv.Api) (*ProductApiClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "productapi", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ProductApiClient: %+v", err)
	}

	return &ProductApiClient{
		Client: client,
	}, nil
}
