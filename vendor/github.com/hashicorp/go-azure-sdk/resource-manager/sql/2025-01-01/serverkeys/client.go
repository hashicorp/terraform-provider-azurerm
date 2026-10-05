package serverkeys

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerKeysClient struct {
	Client *resourcemanager.Client
}

func NewServerKeysClientWithBaseURI(sdkApi sdkEnv.Api) (*ServerKeysClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "serverkeys", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ServerKeysClient: %+v", err)
	}

	return &ServerKeysClient{
		Client: client,
	}, nil
}
