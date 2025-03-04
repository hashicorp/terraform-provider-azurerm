package machines

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MachinesClient struct {
	Client *resourcemanager.Client
}

func NewMachinesClientWithBaseURI(sdkApi sdkEnv.Api) (*MachinesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "machines", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating MachinesClient: %+v", err)
	}

	return &MachinesClient{
		Client: client,
	}, nil
}
