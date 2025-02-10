package virtualharddisks

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualHardDisksClient struct {
	Client *resourcemanager.Client
}

func NewVirtualHardDisksClientWithBaseURI(sdkApi sdkEnv.Api) (*VirtualHardDisksClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "virtualharddisks", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating VirtualHardDisksClient: %+v", err)
	}

	return &VirtualHardDisksClient{
		Client: client,
	}, nil
}
