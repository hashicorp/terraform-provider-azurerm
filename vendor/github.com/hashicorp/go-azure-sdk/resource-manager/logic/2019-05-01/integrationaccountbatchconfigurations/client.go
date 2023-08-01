package integrationaccountbatchconfigurations

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationAccountBatchConfigurationsClient struct {
	Client *resourcemanager.Client
}

func NewIntegrationAccountBatchConfigurationsClientWithBaseURI(sdkApi sdkEnv.Api) (*IntegrationAccountBatchConfigurationsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "integrationaccountbatchconfigurations", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating IntegrationAccountBatchConfigurationsClient: %+v", err)
	}

	return &IntegrationAccountBatchConfigurationsClient{
		Client: client,
	}, nil
}
