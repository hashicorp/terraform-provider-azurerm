package applicationgroup

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGroupClient struct {
	Client *resourcemanager.Client
}

func NewApplicationGroupClientWithBaseURI(api environments.Api) (*ApplicationGroupClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "applicationgroup", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ApplicationGroupClient: %+v", err)
	}

	return &ApplicationGroupClient{
		Client: client,
	}, nil
}
