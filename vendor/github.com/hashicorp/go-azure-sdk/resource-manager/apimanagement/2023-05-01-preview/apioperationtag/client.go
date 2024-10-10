package apioperationtag

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiOperationTagClient struct {
	Client *resourcemanager.Client
}

func NewApiOperationTagClientWithBaseURI(sdkApi sdkEnv.Api) (*ApiOperationTagClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "apioperationtag", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ApiOperationTagClient: %+v", err)
	}

	return &ApiOperationTagClient{
		Client: client,
	}, nil
}
