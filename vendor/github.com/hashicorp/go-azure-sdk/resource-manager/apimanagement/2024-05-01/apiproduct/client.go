package apiproduct

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiProductClient struct {
	Client *resourcemanager.Client
}

func NewApiProductClientWithBaseURI(sdkApi sdkEnv.Api) (*ApiProductClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "apiproduct", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ApiProductClient: %+v", err)
	}

	return &ApiProductClient{
		Client: client,
	}, nil
}
