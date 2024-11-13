package importpipelines

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImportPipelinesClient struct {
	Client *resourcemanager.Client
}

func NewImportPipelinesClientWithBaseURI(sdkApi sdkEnv.Api) (*ImportPipelinesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "importpipelines", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ImportPipelinesClient: %+v", err)
	}

	return &ImportPipelinesClient{
		Client: client,
	}, nil
}
