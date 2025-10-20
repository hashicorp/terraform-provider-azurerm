package dataflow

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataflowClient struct {
	Client *resourcemanager.Client
}

func NewDataflowClientWithBaseURI(sdkApi sdkEnv.Api) (*DataflowClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "dataflow", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DataflowClient: %+v", err)
	}

	return &DataflowClient{
		Client: client,
	}, nil
}
