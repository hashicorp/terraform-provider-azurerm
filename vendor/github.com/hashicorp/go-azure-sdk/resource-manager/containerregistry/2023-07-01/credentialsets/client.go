package credentialsets

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CredentialSetsClient struct {
	Client *resourcemanager.Client
}

func NewCredentialSetsClientWithBaseURI(sdkApi sdkEnv.Api) (*CredentialSetsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "credentialsets", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating CredentialSetsClient: %+v", err)
	}

	return &CredentialSetsClient{
		Client: client,
	}, nil
}
