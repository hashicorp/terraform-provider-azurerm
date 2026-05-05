package deletedvaults

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeletedVaultsClient struct {
	Client *resourcemanager.Client
}

func NewDeletedVaultsClientWithBaseURI(sdkApi sdkEnv.Api) (*DeletedVaultsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "deletedvaults", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DeletedVaultsClient: %+v", err)
	}

	return &DeletedVaultsClient{
		Client: client,
	}, nil
}
