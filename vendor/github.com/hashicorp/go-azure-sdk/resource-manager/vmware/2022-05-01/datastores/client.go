package datastores

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataStoresClient struct {
	Client *resourcemanager.Client
}

func NewDataStoresClientWithBaseURI(sdkApi sdkEnv.Api) (*DataStoresClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "datastores", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DataStoresClient: %+v", err)
	}

	return &DataStoresClient{
		Client: client,
	}, nil
}
