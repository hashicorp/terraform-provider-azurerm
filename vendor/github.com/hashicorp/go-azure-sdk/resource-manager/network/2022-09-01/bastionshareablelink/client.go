package bastionshareablelink

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BastionShareableLinkClient struct {
	Client *resourcemanager.Client
}

func NewBastionShareableLinkClientWithBaseURI(api environments.Api) (*BastionShareableLinkClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "bastionshareablelink", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating BastionShareableLinkClient: %+v", err)
	}

	return &BastionShareableLinkClient{
		Client: client,
	}, nil
}
