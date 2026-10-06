package managedenvironmentsstorages

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedEnvironmentsStoragesClient struct {
	Client *resourcemanager.Client
}

func NewManagedEnvironmentsStoragesClientWithBaseURI(sdkApi sdkEnv.Api) (*ManagedEnvironmentsStoragesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "managedenvironmentsstorages", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ManagedEnvironmentsStoragesClient: %+v", err)
	}

	return &ManagedEnvironmentsStoragesClient{
		Client: client,
	}, nil
}
