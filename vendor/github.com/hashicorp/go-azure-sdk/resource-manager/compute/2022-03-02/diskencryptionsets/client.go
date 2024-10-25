package diskencryptionsets

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiskEncryptionSetsClient struct {
	Client *resourcemanager.Client
}

func NewDiskEncryptionSetsClientWithBaseURI(sdkApi sdkEnv.Api) (*DiskEncryptionSetsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "diskencryptionsets", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DiskEncryptionSetsClient: %+v", err)
	}

	return &DiskEncryptionSetsClient{
		Client: client,
	}, nil
}
