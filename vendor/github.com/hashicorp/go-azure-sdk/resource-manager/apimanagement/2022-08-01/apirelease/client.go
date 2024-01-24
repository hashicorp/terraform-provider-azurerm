package apirelease

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiReleaseClient struct {
	Client *resourcemanager.Client
}

func NewApiReleaseClientWithBaseURI(sdkApi sdkEnv.Api) (*ApiReleaseClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "apirelease", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ApiReleaseClient: %+v", err)
	}

	return &ApiReleaseClient{
		Client: client,
	}, nil
}
