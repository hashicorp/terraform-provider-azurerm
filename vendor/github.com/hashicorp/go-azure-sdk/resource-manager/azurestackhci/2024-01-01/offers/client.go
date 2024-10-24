package offers

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OffersClient struct {
	Client *resourcemanager.Client
}

func NewOffersClientWithBaseURI(sdkApi sdkEnv.Api) (*OffersClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "offers", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating OffersClient: %+v", err)
	}

	return &OffersClient{
		Client: client,
	}, nil
}
