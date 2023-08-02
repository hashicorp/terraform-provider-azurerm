package topics

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TopicsClient struct {
	Client *resourcemanager.Client
}

func NewTopicsClientWithBaseURI(sdkApi sdkEnv.Api) (*TopicsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "topics", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating TopicsClient: %+v", err)
	}

	return &TopicsClient{
		Client: client,
	}, nil
}
