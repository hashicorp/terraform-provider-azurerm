package federatedidentitycredentials

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FederatedIdentityCredentialsClient struct {
	Client *resourcemanager.Client
}

func NewFederatedIdentityCredentialsClientWithBaseURI(sdkApi sdkEnv.Api) (*FederatedIdentityCredentialsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "federatedidentitycredentials", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating FederatedIdentityCredentialsClient: %+v", err)
	}

	return &FederatedIdentityCredentialsClient{
		Client: client,
	}, nil
}
