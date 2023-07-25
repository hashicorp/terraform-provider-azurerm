package staticmembers

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StaticMembersClient struct {
	Client *resourcemanager.Client
}

func NewStaticMembersClientWithBaseURI(api environments.Api) (*StaticMembersClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "staticmembers", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating StaticMembersClient: %+v", err)
	}

	return &StaticMembersClient{
		Client: client,
	}, nil
}
