package machinenetworkprofile

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MachineNetworkProfileClient struct {
	Client *resourcemanager.Client
}

func NewMachineNetworkProfileClientWithBaseURI(sdkApi sdkEnv.Api) (*MachineNetworkProfileClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "machinenetworkprofile", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating MachineNetworkProfileClient: %+v", err)
	}

	return &MachineNetworkProfileClient{
		Client: client,
	}, nil
}
