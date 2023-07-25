package factories

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FactoriesClient struct {
	Client *resourcemanager.Client
}

func NewFactoriesClientWithBaseURI(api environments.Api) (*FactoriesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "factories", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating FactoriesClient: %+v", err)
	}

	return &FactoriesClient{
		Client: client,
	}, nil
}
