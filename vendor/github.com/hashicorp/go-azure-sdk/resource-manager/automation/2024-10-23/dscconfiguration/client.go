package dscconfiguration

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DscConfigurationClient struct {
	Client *resourcemanager.Client
}

func NewDscConfigurationClientWithBaseURI(sdkApi sdkEnv.Api) (*DscConfigurationClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "dscconfiguration", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DscConfigurationClient: %+v", err)
	}

	return &DscConfigurationClient{
		Client: client,
	}, nil
}
