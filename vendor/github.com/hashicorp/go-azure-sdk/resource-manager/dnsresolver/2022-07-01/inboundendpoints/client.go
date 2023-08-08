package inboundendpoints

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InboundEndpointsClient struct {
	Client *resourcemanager.Client
}

func NewInboundEndpointsClientWithBaseURI(sdkApi sdkEnv.Api) (*InboundEndpointsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "inboundendpoints", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating InboundEndpointsClient: %+v", err)
	}

	return &InboundEndpointsClient{
		Client: client,
	}, nil
}
