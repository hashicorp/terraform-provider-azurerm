package dscnodeconfiguration

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DscNodeConfigurationClient struct {
	Client *resourcemanager.Client
}

func NewDscNodeConfigurationClientWithBaseURI(sdkApi sdkEnv.Api) (*DscNodeConfigurationClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "dscnodeconfiguration", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DscNodeConfigurationClient: %+v", err)
	}

	return &DscNodeConfigurationClient{
		Client: client,
	}, nil
}
