package tenants

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TenantsClient struct {
	Client *resourcemanager.Client
}

func NewTenantsClientWithBaseURI(sdkApi sdkEnv.Api) (*TenantsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "tenants", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating TenantsClient: %+v", err)
	}

	return &TenantsClient{
		Client: client,
	}, nil
}
