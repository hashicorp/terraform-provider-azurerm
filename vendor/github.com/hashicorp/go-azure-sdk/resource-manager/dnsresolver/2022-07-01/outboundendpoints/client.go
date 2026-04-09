package outboundendpoints

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OutboundEndpointsClient struct {
	Client *resourcemanager.Client
}

func NewOutboundEndpointsClientWithBaseURI(sdkApi sdkEnv.Api) (*OutboundEndpointsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "outboundendpoints", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating OutboundEndpointsClient: %+v", err)
	}

	return &OutboundEndpointsClient{
		Client: client,
	}, nil
}
