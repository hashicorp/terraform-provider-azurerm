package databasesystemshaperesources

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabaseSystemShapeResourcesClient struct {
	Client *resourcemanager.Client
}

func NewDatabaseSystemShapeResourcesClientWithBaseURI(sdkApi sdkEnv.Api) (*DatabaseSystemShapeResourcesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "databasesystemshaperesources", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DatabaseSystemShapeResourcesClient: %+v", err)
	}

	return &DatabaseSystemShapeResourcesClient{
		Client: client,
	}, nil
}
