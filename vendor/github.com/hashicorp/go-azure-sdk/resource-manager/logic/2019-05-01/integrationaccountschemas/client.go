package integrationaccountschemas

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationAccountSchemasClient struct {
	Client *resourcemanager.Client
}

func NewIntegrationAccountSchemasClientWithBaseURI(api environments.Api) (*IntegrationAccountSchemasClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "integrationaccountschemas", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating IntegrationAccountSchemasClient: %+v", err)
	}

	return &IntegrationAccountSchemasClient{
		Client: client,
	}, nil
}
