package account

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountClient struct {
	Client *resourcemanager.Client
}

func NewAccountClientWithBaseURI(sdkApi sdkEnv.Api) (*AccountClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "account", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AccountClient: %+v", err)
	}

	return &AccountClient{
		Client: client,
	}, nil
}
