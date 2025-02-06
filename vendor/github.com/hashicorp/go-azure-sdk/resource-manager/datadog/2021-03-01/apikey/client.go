package apikey

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiKeyClient struct {
	Client *resourcemanager.Client
}

func NewApiKeyClientWithBaseURI(sdkApi sdkEnv.Api) (*ApiKeyClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "apikey", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ApiKeyClient: %+v", err)
	}

	return &ApiKeyClient{
		Client: client,
	}, nil
}
