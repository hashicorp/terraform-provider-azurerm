package queues

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueuesClient struct {
	Client *resourcemanager.Client
}

func NewQueuesClientWithBaseURI(api environments.Api) (*QueuesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "queues", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating QueuesClient: %+v", err)
	}

	return &QueuesClient{
		Client: client,
	}, nil
}
