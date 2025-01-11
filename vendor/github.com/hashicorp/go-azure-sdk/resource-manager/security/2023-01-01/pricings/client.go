package pricings

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PricingsClient struct {
	Client *resourcemanager.Client
}

func NewPricingsClientWithBaseURI(sdkApi sdkEnv.Api) (*PricingsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "pricings", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PricingsClient: %+v", err)
	}

	return &PricingsClient{
		Client: client,
	}, nil
}
