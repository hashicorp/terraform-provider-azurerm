package exascaledbstoragevaults

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExascaleDbStorageVaultsClient struct {
	Client *resourcemanager.Client
}

func NewExascaleDbStorageVaultsClientWithBaseURI(sdkApi sdkEnv.Api) (*ExascaleDbStorageVaultsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "exascaledbstoragevaults", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ExascaleDbStorageVaultsClient: %+v", err)
	}

	return &ExascaleDbStorageVaultsClient{
		Client: client,
	}, nil
}
