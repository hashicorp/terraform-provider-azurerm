package tableservice

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TableServiceClient struct {
	Client *resourcemanager.Client
}

func NewTableServiceClientWithBaseURI(sdkApi sdkEnv.Api) (*TableServiceClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "tableservice", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating TableServiceClient: %+v", err)
	}

	return &TableServiceClient{
		Client: client,
	}, nil
}
