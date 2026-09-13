package blobinventorypolicies

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BlobInventoryPoliciesClient struct {
	Client *resourcemanager.Client
}

func NewBlobInventoryPoliciesClientWithBaseURI(sdkApi sdkEnv.Api) (*BlobInventoryPoliciesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "blobinventorypolicies", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating BlobInventoryPoliciesClient: %+v", err)
	}

	return &BlobInventoryPoliciesClient{
		Client: client,
	}, nil
}
