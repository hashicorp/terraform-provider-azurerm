package dataflowendpoint

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataflowEndpointClient struct {
	Client *resourcemanager.Client
}

func NewDataflowEndpointClientWithBaseURI(sdkApi sdkEnv.Api) (*DataflowEndpointClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "dataflowendpoint", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DataflowEndpointClient: %+v", err)
	}

	return &DataflowEndpointClient{
		Client: client,
	}, nil
}
