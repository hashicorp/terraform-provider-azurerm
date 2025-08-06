package managedinstanceadministrators

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedInstanceAdministratorsClient struct {
	Client *resourcemanager.Client
}

func NewManagedInstanceAdministratorsClientWithBaseURI(sdkApi sdkEnv.Api) (*ManagedInstanceAdministratorsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "managedinstanceadministrators", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ManagedInstanceAdministratorsClient: %+v", err)
	}

	return &ManagedInstanceAdministratorsClient{
		Client: client,
	}, nil
}
