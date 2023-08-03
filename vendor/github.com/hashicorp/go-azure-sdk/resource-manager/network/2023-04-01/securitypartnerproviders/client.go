package securitypartnerproviders

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityPartnerProvidersClient struct {
	Client *resourcemanager.Client
}

func NewSecurityPartnerProvidersClientWithBaseURI(sdkApi sdkEnv.Api) (*SecurityPartnerProvidersClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "securitypartnerproviders", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SecurityPartnerProvidersClient: %+v", err)
	}

	return &SecurityPartnerProvidersClient{
		Client: client,
	}, nil
}
