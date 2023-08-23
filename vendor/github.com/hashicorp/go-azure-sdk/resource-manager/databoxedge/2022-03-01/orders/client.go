package orders

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OrdersClient struct {
	Client *resourcemanager.Client
}

func NewOrdersClientWithBaseURI(sdkApi sdkEnv.Api) (*OrdersClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "orders", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating OrdersClient: %+v", err)
	}

	return &OrdersClient{
		Client: client,
	}, nil
}
