package applicationgateways

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewaysClient struct {
	Client *resourcemanager.Client
}

func NewApplicationGatewaysClientWithBaseURI(sdkApi sdkEnv.Api) (*ApplicationGatewaysClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "applicationgateways", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ApplicationGatewaysClient: %+v", err)
	}

	return &ApplicationGatewaysClient{
		Client: client,
	}, nil
}
