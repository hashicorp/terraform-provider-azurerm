package apimanagementgatewayskus

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiManagementGatewaySkusClient struct {
	Client *resourcemanager.Client
}

func NewApiManagementGatewaySkusClientWithBaseURI(sdkApi sdkEnv.Api) (*ApiManagementGatewaySkusClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "apimanagementgatewayskus", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ApiManagementGatewaySkusClient: %+v", err)
	}

	return &ApiManagementGatewaySkusClient{
		Client: client,
	}, nil
}
