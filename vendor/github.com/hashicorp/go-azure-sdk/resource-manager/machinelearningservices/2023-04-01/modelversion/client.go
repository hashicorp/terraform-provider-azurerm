package modelversion

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ModelVersionClient struct {
	Client *resourcemanager.Client
}

func NewModelVersionClientWithBaseURI(api environments.Api) (*ModelVersionClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "modelversion", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ModelVersionClient: %+v", err)
	}

	return &ModelVersionClient{
		Client: client,
	}, nil
}
