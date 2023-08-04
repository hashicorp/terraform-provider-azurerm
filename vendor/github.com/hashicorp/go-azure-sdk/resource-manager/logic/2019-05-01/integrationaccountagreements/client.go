package integrationaccountagreements

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationAccountAgreementsClient struct {
	Client *resourcemanager.Client
}

func NewIntegrationAccountAgreementsClientWithBaseURI(sdkApi sdkEnv.Api) (*IntegrationAccountAgreementsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "integrationaccountagreements", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating IntegrationAccountAgreementsClient: %+v", err)
	}

	return &IntegrationAccountAgreementsClient{
		Client: client,
	}, nil
}
