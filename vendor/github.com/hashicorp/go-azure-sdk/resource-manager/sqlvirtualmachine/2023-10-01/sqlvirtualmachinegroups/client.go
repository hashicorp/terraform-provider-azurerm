package sqlvirtualmachinegroups

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SqlVirtualMachineGroupsClient struct {
	Client *resourcemanager.Client
}

func NewSqlVirtualMachineGroupsClientWithBaseURI(sdkApi sdkEnv.Api) (*SqlVirtualMachineGroupsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "sqlvirtualmachinegroups", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SqlVirtualMachineGroupsClient: %+v", err)
	}

	return &SqlVirtualMachineGroupsClient{
		Client: client,
	}, nil
}
