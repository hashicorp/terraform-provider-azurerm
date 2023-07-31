package serverstart

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerStartClient struct {
	Client *resourcemanager.Client
}

func NewServerStartClientWithBaseURI(api environments.Api) (*ServerStartClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "serverstart", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ServerStartClient: %+v", err)
	}

	return &ServerStartClient{
		Client: client,
	}, nil
}
