package channel

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ChannelClient struct {
	Client *resourcemanager.Client
}

func NewChannelClientWithBaseURI(sdkApi sdkEnv.Api) (*ChannelClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "channel", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ChannelClient: %+v", err)
	}

	return &ChannelClient{
		Client: client,
	}, nil
}
