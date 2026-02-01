package apimanagementservice

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiManagementServiceClient struct {
	Client *resourcemanager.Client
}

func NewApiManagementServiceClientWithBaseURI(sdkApi sdkEnv.Api) (*ApiManagementServiceClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "apimanagementservice", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ApiManagementServiceClient: %+v", err)
	}

	return &ApiManagementServiceClient{
		Client: client,
	}, nil
}
