package querypackqueries

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueryPackQueriesClient struct {
	Client *resourcemanager.Client
}

func NewQueryPackQueriesClientWithBaseURI(sdkApi sdkEnv.Api) (*QueryPackQueriesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "querypackqueries", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating QueryPackQueriesClient: %+v", err)
	}

	return &QueryPackQueriesClient{
		Client: client,
	}, nil
}
