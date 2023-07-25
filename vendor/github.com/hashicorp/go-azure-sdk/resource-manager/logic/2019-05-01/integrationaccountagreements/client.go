package integrationaccountagreements

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationAccountAgreementsClient struct {
	Client *resourcemanager.Client
}

func NewIntegrationAccountAgreementsClientWithBaseURI(api environments.Api) (*IntegrationAccountAgreementsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "integrationaccountagreements", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating IntegrationAccountAgreementsClient: %+v", err)
	}

	return &IntegrationAccountAgreementsClient{
		Client: client,
	}, nil
}
