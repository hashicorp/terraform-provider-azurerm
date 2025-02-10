package views

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ViewsClient struct {
	Client *resourcemanager.Client
}

func NewViewsClientWithBaseURI(sdkApi sdkEnv.Api) (*ViewsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "views", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ViewsClient: %+v", err)
	}

	return &ViewsClient{
		Client: client,
	}, nil
}
