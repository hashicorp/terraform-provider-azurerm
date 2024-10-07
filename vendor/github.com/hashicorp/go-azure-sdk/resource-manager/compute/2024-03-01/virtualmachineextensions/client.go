package virtualmachineextensions

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineExtensionsClient struct {
	Client *resourcemanager.Client
}

func NewVirtualMachineExtensionsClientWithBaseURI(sdkApi sdkEnv.Api) (*VirtualMachineExtensionsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "virtualmachineextensions", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating VirtualMachineExtensionsClient: %+v", err)
	}

	return &VirtualMachineExtensionsClient{
		Client: client,
	}, nil
}
