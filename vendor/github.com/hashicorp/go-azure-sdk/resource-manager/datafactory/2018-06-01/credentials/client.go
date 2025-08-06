package credentials

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CredentialsClient struct {
	Client *resourcemanager.Client
}

func NewCredentialsClientWithBaseURI(sdkApi sdkEnv.Api) (*CredentialsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "credentials", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating CredentialsClient: %+v", err)
	}

	return &CredentialsClient{
		Client: client,
	}, nil
}
