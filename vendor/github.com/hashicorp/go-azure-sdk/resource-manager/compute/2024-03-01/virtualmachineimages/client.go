package virtualmachineimages

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineImagesClient struct {
	Client *resourcemanager.Client
}

func NewVirtualMachineImagesClientWithBaseURI(sdkApi sdkEnv.Api) (*VirtualMachineImagesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "virtualmachineimages", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating VirtualMachineImagesClient: %+v", err)
	}

	return &VirtualMachineImagesClient{
		Client: client,
	}, nil
}
