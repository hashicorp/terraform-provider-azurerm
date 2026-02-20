package mongorbacs

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongorbacsClient struct {
	Client *resourcemanager.Client
}

func NewMongorbacsClientWithBaseURI(sdkApi sdkEnv.Api) (*MongorbacsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "mongorbacs", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating MongorbacsClient: %+v", err)
	}

	return &MongorbacsClient{
		Client: client,
	}, nil
}
