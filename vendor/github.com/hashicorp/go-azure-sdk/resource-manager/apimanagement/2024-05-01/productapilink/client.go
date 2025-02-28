package productapilink

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProductApiLinkClient struct {
	Client *resourcemanager.Client
}

func NewProductApiLinkClientWithBaseURI(sdkApi sdkEnv.Api) (*ProductApiLinkClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "productapilink", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ProductApiLinkClient: %+v", err)
	}

	return &ProductApiLinkClient{
		Client: client,
	}, nil
}
