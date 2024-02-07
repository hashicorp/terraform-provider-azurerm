package subaccount

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubAccountClient struct {
	Client *resourcemanager.Client
}

func NewSubAccountClientWithBaseURI(sdkApi sdkEnv.Api) (*SubAccountClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "subaccount", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SubAccountClient: %+v", err)
	}

	return &SubAccountClient{
		Client: client,
	}, nil
}
