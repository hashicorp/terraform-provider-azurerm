package apioperationsbytag

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiOperationsByTagClient struct {
	Client *resourcemanager.Client
}

func NewApiOperationsByTagClientWithBaseURI(sdkApi sdkEnv.Api) (*ApiOperationsByTagClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "apioperationsbytag", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ApiOperationsByTagClient: %+v", err)
	}

	return &ApiOperationsByTagClient{
		Client: client,
	}, nil
}
