package securityadminconfigurations

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityAdminConfigurationsClient struct {
	Client *resourcemanager.Client
}

func NewSecurityAdminConfigurationsClientWithBaseURI(sdkApi sdkEnv.Api) (*SecurityAdminConfigurationsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "securityadminconfigurations", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SecurityAdminConfigurationsClient: %+v", err)
	}

	return &SecurityAdminConfigurationsClient{
		Client: client,
	}, nil
}
