package virtualmachines

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachinesClient struct {
	Client *resourcemanager.Client
}

func NewVirtualMachinesClientWithBaseURI(sdkApi sdkEnv.Api) (*VirtualMachinesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "virtualmachines", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating VirtualMachinesClient: %+v", err)
	}

	return &VirtualMachinesClient{
		Client: client,
	}, nil
}
