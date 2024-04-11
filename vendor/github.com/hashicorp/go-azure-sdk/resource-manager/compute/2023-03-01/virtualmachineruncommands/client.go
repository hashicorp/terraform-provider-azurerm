package virtualmachineruncommands

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineRunCommandsClient struct {
	Client *resourcemanager.Client
}

func NewVirtualMachineRunCommandsClientWithBaseURI(sdkApi sdkEnv.Api) (*VirtualMachineRunCommandsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "virtualmachineruncommands", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating VirtualMachineRunCommandsClient: %+v", err)
	}

	return &VirtualMachineRunCommandsClient{
		Client: client,
	}, nil
}
