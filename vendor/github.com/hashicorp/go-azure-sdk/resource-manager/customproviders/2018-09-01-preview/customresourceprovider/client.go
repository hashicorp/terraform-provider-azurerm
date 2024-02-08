package customresourceprovider

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomResourceProviderClient struct {
	Client *resourcemanager.Client
}

func NewCustomResourceProviderClientWithBaseURI(sdkApi sdkEnv.Api) (*CustomResourceProviderClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "customresourceprovider", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating CustomResourceProviderClient: %+v", err)
	}

	return &CustomResourceProviderClient{
		Client: client,
	}, nil
}
