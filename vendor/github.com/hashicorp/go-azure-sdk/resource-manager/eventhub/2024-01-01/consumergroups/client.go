package consumergroups

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConsumerGroupsClient struct {
	Client *resourcemanager.Client
}

func NewConsumerGroupsClientWithBaseURI(sdkApi sdkEnv.Api) (*ConsumerGroupsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "consumergroups", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ConsumerGroupsClient: %+v", err)
	}

	return &ConsumerGroupsClient{
		Client: client,
	}, nil
}
