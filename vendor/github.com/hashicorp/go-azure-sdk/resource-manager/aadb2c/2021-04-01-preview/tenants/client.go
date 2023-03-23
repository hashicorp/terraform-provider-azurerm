package tenants

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TenantsClient struct {
	Client *resourcemanager.Client
}

func NewTenantsClientWithBaseURI(api environments.Api) (*TenantsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "tenants", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating TenantsClient: %+v", err)
	}

	return &TenantsClient{
		Client: client,
	}, nil
}
