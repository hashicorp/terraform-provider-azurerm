package virtualmachinescalesetextensions

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineScaleSetExtensionsClient struct {
	Client *resourcemanager.Client
}

func NewVirtualMachineScaleSetExtensionsClientWithBaseURI(sdkApi sdkEnv.Api) (*VirtualMachineScaleSetExtensionsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "virtualmachinescalesetextensions", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating VirtualMachineScaleSetExtensionsClient: %+v", err)
	}

	return &VirtualMachineScaleSetExtensionsClient{
		Client: client,
	}, nil
}
