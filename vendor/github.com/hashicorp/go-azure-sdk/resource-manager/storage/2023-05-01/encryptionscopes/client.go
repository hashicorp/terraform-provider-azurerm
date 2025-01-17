package encryptionscopes

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EncryptionScopesClient struct {
	Client *resourcemanager.Client
}

func NewEncryptionScopesClientWithBaseURI(sdkApi sdkEnv.Api) (*EncryptionScopesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "encryptionscopes", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating EncryptionScopesClient: %+v", err)
	}

	return &EncryptionScopesClient{
		Client: client,
	}, nil
}
