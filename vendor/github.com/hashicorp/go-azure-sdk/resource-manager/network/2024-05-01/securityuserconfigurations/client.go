package securityuserconfigurations

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityUserConfigurationsClient struct {
	Client *resourcemanager.Client
}

func NewSecurityUserConfigurationsClientWithBaseURI(sdkApi sdkEnv.Api) (*SecurityUserConfigurationsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "securityuserconfigurations", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SecurityUserConfigurationsClient: %+v", err)
	}

	return &SecurityUserConfigurationsClient{
		Client: client,
	}, nil
}
