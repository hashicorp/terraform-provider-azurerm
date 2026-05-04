package serverfailover

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerFailoverClient struct {
	Client *resourcemanager.Client
}

func NewServerFailoverClientWithBaseURI(sdkApi sdkEnv.Api) (*ServerFailoverClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "serverfailover", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ServerFailoverClient: %+v", err)
	}

	return &ServerFailoverClient{
		Client: client,
	}, nil
}
