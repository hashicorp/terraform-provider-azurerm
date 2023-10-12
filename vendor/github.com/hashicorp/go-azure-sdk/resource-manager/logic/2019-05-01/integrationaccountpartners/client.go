package integrationaccountpartners

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationAccountPartnersClient struct {
	Client *resourcemanager.Client
}

func NewIntegrationAccountPartnersClientWithBaseURI(sdkApi sdkEnv.Api) (*IntegrationAccountPartnersClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "integrationaccountpartners", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating IntegrationAccountPartnersClient: %+v", err)
	}

	return &IntegrationAccountPartnersClient{
		Client: client,
	}, nil
}
