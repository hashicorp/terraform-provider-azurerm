package expressroutelagauthorizations

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteLagAuthorizationsClient struct {
	Client *resourcemanager.Client
}

func NewExpressRouteLagAuthorizationsClientWithBaseURI(sdkApi sdkEnv.Api) (*ExpressRouteLagAuthorizationsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "expressroutelagauthorizations", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ExpressRouteLagAuthorizationsClient: %+v", err)
	}

	return &ExpressRouteLagAuthorizationsClient{
		Client: client,
	}, nil
}
