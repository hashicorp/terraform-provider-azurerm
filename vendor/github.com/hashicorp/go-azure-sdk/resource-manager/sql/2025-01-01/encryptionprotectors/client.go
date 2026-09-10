package encryptionprotectors

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EncryptionProtectorsClient struct {
	Client *resourcemanager.Client
}

func NewEncryptionProtectorsClientWithBaseURI(sdkApi sdkEnv.Api) (*EncryptionProtectorsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "encryptionprotectors", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating EncryptionProtectorsClient: %+v", err)
	}

	return &EncryptionProtectorsClient{
		Client: client,
	}, nil
}
