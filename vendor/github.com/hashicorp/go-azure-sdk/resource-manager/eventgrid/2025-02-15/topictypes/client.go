package topictypes

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TopicTypesClient struct {
	Client *resourcemanager.Client
}

func NewTopicTypesClientWithBaseURI(sdkApi sdkEnv.Api) (*TopicTypesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "topictypes", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating TopicTypesClient: %+v", err)
	}

	return &TopicTypesClient{
		Client: client,
	}, nil
}
