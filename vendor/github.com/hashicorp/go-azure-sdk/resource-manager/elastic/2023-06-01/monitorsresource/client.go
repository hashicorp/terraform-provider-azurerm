package monitorsresource

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MonitorsResourceClient struct {
	Client *resourcemanager.Client
}

func NewMonitorsResourceClientWithBaseURI(api environments.Api) (*MonitorsResourceClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "monitorsresource", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating MonitorsResourceClient: %+v", err)
	}

	return &MonitorsResourceClient{
		Client: client,
	}, nil
}
