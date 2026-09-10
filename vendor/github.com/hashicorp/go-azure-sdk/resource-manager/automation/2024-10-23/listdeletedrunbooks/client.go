package listdeletedrunbooks

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListDeletedRunbooksClient struct {
	Client *resourcemanager.Client
}

func NewListDeletedRunbooksClientWithBaseURI(sdkApi sdkEnv.Api) (*ListDeletedRunbooksClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "listdeletedrunbooks", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ListDeletedRunbooksClient: %+v", err)
	}

	return &ListDeletedRunbooksClient{
		Client: client,
	}, nil
}
