package sharedprivatelinkresources

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SharedPrivateLinkResourcesClient struct {
	Client *resourcemanager.Client
}

func NewSharedPrivateLinkResourcesClientWithBaseURI(sdkApi sdkEnv.Api) (*SharedPrivateLinkResourcesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "sharedprivatelinkresources", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SharedPrivateLinkResourcesClient: %+v", err)
	}

	return &SharedPrivateLinkResourcesClient{
		Client: client,
	}, nil
}
