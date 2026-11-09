package expressroutelagmembers

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteLagMembersClient struct {
	Client *resourcemanager.Client
}

func NewExpressRouteLagMembersClientWithBaseURI(sdkApi sdkEnv.Api) (*ExpressRouteLagMembersClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "expressroutelagmembers", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ExpressRouteLagMembersClient: %+v", err)
	}

	return &ExpressRouteLagMembersClient{
		Client: client,
	}, nil
}
