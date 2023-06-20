package desktop

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DesktopClient struct {
	Client *resourcemanager.Client
}

func NewDesktopClientWithBaseURI(api environments.Api) (*DesktopClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "desktop", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DesktopClient: %+v", err)
	}

	return &DesktopClient{
		Client: client,
	}, nil
}
