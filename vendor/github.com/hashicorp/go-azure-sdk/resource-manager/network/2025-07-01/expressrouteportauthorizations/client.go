package expressrouteportauthorizations

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRoutePortAuthorizationsClient struct {
	Client *resourcemanager.Client
}

func NewExpressRoutePortAuthorizationsClientWithBaseURI(sdkApi sdkEnv.Api) (*ExpressRoutePortAuthorizationsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "expressrouteportauthorizations", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ExpressRoutePortAuthorizationsClient: %+v", err)
	}

	return &ExpressRoutePortAuthorizationsClient{
		Client: client,
	}, nil
}
