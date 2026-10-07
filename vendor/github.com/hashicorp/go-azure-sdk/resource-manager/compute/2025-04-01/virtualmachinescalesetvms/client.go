package virtualmachinescalesetvms

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineScaleSetVMsClient struct {
	Client *resourcemanager.Client
}

func NewVirtualMachineScaleSetVMsClientWithBaseURI(sdkApi sdkEnv.Api) (*VirtualMachineScaleSetVMsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "virtualmachinescalesetvms", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating VirtualMachineScaleSetVMsClient: %+v", err)
	}

	return &VirtualMachineScaleSetVMsClient{
		Client: client,
	}, nil
}
