package dbservers

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DbServersClient struct {
	Client *resourcemanager.Client
}

func NewDbServersClientWithBaseURI(sdkApi sdkEnv.Api) (*DbServersClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "dbservers", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DbServersClient: %+v", err)
	}

	return &DbServersClient{
		Client: client,
	}, nil
}
