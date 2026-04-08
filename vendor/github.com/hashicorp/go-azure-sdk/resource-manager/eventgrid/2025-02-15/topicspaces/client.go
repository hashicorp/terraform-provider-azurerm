package topicspaces

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TopicSpacesClient struct {
	Client *resourcemanager.Client
}

func NewTopicSpacesClientWithBaseURI(sdkApi sdkEnv.Api) (*TopicSpacesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "topicspaces", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating TopicSpacesClient: %+v", err)
	}

	return &TopicSpacesClient{
		Client: client,
	}, nil
}
