package productsbytag

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProductsByTagClient struct {
	Client *resourcemanager.Client
}

func NewProductsByTagClientWithBaseURI(sdkApi sdkEnv.Api) (*ProductsByTagClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "productsbytag", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ProductsByTagClient: %+v", err)
	}

	return &ProductsByTagClient{
		Client: client,
	}, nil
}
