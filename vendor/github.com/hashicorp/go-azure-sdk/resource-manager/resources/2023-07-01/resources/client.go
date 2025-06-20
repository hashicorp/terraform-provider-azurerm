package resources

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourcesClient struct {
	Client *resourcemanager.Client
}

func NewResourcesClientWithBaseURI(sdkApi sdkEnv.Api) (*ResourcesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "resources", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ResourcesClient: %+v", err)
	}

	return &ResourcesClient{
		Client: client,
	}, nil
}
