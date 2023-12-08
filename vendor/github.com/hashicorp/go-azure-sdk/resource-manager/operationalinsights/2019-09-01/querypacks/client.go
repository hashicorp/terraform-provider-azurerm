package querypacks

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueryPacksClient struct {
	Client *resourcemanager.Client
}

func NewQueryPacksClientWithBaseURI(sdkApi sdkEnv.Api) (*QueryPacksClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "querypacks", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating QueryPacksClient: %+v", err)
	}

	return &QueryPacksClient{
		Client: client,
	}, nil
}
