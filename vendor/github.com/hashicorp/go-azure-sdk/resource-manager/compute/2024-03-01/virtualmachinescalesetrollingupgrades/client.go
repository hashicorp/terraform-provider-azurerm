package virtualmachinescalesetrollingupgrades

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineScaleSetRollingUpgradesClient struct {
	Client *resourcemanager.Client
}

func NewVirtualMachineScaleSetRollingUpgradesClientWithBaseURI(sdkApi sdkEnv.Api) (*VirtualMachineScaleSetRollingUpgradesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "virtualmachinescalesetrollingupgrades", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating VirtualMachineScaleSetRollingUpgradesClient: %+v", err)
	}

	return &VirtualMachineScaleSetRollingUpgradesClient{
		Client: client,
	}, nil
}
