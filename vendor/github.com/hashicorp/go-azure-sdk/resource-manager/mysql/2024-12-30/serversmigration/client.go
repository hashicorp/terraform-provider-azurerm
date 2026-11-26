package serversmigration

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServersMigrationClient struct {
	Client *resourcemanager.Client
}

func NewServersMigrationClientWithBaseURI(sdkApi sdkEnv.Api) (*ServersMigrationClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "serversmigration", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ServersMigrationClient: %+v", err)
	}

	return &ServersMigrationClient{
		Client: client,
	}, nil
}
