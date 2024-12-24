package graphqueries

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GraphqueriesClient struct {
	Client *resourcemanager.Client
}

func NewGraphqueriesClientWithBaseURI(sdkApi sdkEnv.Api) (*GraphqueriesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "graphqueries", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating GraphqueriesClient: %+v", err)
	}

	return &GraphqueriesClient{
		Client: client,
	}, nil
}
