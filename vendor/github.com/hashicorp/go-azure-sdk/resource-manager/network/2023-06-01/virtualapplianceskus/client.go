package virtualapplianceskus

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualApplianceSkusClient struct {
	Client *resourcemanager.Client
}

func NewVirtualApplianceSkusClientWithBaseURI(sdkApi sdkEnv.Api) (*VirtualApplianceSkusClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "virtualapplianceskus", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating VirtualApplianceSkusClient: %+v", err)
	}

	return &VirtualApplianceSkusClient{
		Client: client,
	}, nil
}
