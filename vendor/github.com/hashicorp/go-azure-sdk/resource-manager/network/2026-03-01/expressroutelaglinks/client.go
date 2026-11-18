package expressroutelaglinks

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteLagLinksClient struct {
	Client *resourcemanager.Client
}

func NewExpressRouteLagLinksClientWithBaseURI(sdkApi sdkEnv.Api) (*ExpressRouteLagLinksClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "expressroutelaglinks", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ExpressRouteLagLinksClient: %+v", err)
	}

	return &ExpressRouteLagLinksClient{
		Client: client,
	}, nil
}
