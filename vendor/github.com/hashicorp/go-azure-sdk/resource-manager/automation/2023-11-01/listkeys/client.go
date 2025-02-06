package listkeys

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListKeysClient struct {
	Client *resourcemanager.Client
}

func NewListKeysClientWithBaseURI(sdkApi sdkEnv.Api) (*ListKeysClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "listkeys", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ListKeysClient: %+v", err)
	}

	return &ListKeysClient{
		Client: client,
	}, nil
}
