package querykeys

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueryKeysClient struct {
	Client *resourcemanager.Client
}

func NewQueryKeysClientWithBaseURI(sdkApi sdkEnv.Api) (*QueryKeysClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "querykeys", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating QueryKeysClient: %+v", err)
	}

	return &QueryKeysClient{
		Client: client,
	}, nil
}
