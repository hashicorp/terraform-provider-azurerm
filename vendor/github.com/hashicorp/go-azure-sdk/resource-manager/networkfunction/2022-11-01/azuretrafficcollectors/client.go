package azuretrafficcollectors

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureTrafficCollectorsClient struct {
	Client *resourcemanager.Client
}

func NewAzureTrafficCollectorsClientWithBaseURI(sdkApi sdkEnv.Api) (*AzureTrafficCollectorsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "azuretrafficcollectors", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AzureTrafficCollectorsClient: %+v", err)
	}

	return &AzureTrafficCollectorsClient{
		Client: client,
	}, nil
}
