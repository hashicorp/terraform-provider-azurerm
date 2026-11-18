package edgeactionexecutionfilters

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EdgeActionExecutionFiltersClient struct {
	Client *resourcemanager.Client
}

func NewEdgeActionExecutionFiltersClientWithBaseURI(sdkApi sdkEnv.Api) (*EdgeActionExecutionFiltersClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "edgeactionexecutionfilters", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating EdgeActionExecutionFiltersClient: %+v", err)
	}

	return &EdgeActionExecutionFiltersClient{
		Client: client,
	}, nil
}
