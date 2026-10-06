package channels

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ChannelsClient struct {
	Client *resourcemanager.Client
}

func NewChannelsClientWithBaseURI(sdkApi sdkEnv.Api) (*ChannelsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "channels", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ChannelsClient: %+v", err)
	}

	return &ChannelsClient{
		Client: client,
	}, nil
}
