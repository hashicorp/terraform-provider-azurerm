package resourceguardresources

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceGuardResourcesClient struct {
	Client *resourcemanager.Client
}

func NewResourceGuardResourcesClientWithBaseURI(sdkApi sdkEnv.Api) (*ResourceGuardResourcesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "resourceguardresources", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ResourceGuardResourcesClient: %+v", err)
	}

	return &ResourceGuardResourcesClient{
		Client: client,
	}, nil
}
