package databaseeditions

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabaseEditionsClient struct {
	Client *resourcemanager.Client
}

func NewDatabaseEditionsClientWithBaseURI(sdkApi sdkEnv.Api) (*DatabaseEditionsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "databaseeditions", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DatabaseEditionsClient: %+v", err)
	}

	return &DatabaseEditionsClient{
		Client: client,
	}, nil
}
