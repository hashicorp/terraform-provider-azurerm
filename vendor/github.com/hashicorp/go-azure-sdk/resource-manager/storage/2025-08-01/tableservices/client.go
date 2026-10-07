package tableservices

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TableServicesClient struct {
	Client *resourcemanager.Client
}

func NewTableServicesClientWithBaseURI(sdkApi sdkEnv.Api) (*TableServicesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "tableservices", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating TableServicesClient: %+v", err)
	}

	return &TableServicesClient{
		Client: client,
	}, nil
}
