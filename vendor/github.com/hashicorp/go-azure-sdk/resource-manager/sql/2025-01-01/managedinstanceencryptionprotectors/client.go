package managedinstanceencryptionprotectors

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedInstanceEncryptionProtectorsClient struct {
	Client *resourcemanager.Client
}

func NewManagedInstanceEncryptionProtectorsClientWithBaseURI(sdkApi sdkEnv.Api) (*ManagedInstanceEncryptionProtectorsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "managedinstanceencryptionprotectors", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ManagedInstanceEncryptionProtectorsClient: %+v", err)
	}

	return &ManagedInstanceEncryptionProtectorsClient{
		Client: client,
	}, nil
}
