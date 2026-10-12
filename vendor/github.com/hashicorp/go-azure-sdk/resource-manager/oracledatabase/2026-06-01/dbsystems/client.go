package dbsystems

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DbSystemsClient struct {
	Client *resourcemanager.Client
}

func NewDbSystemsClientWithBaseURI(sdkApi sdkEnv.Api) (*DbSystemsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "dbsystems", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DbSystemsClient: %+v", err)
	}

	return &DbSystemsClient{
		Client: client,
	}, nil
}
