package datastore

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatastoreClient struct {
	Client *resourcemanager.Client
}

func NewDatastoreClientWithBaseURI(sdkApi sdkEnv.Api) (*DatastoreClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "datastore", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DatastoreClient: %+v", err)
	}

	return &DatastoreClient{
		Client: client,
	}, nil
}
