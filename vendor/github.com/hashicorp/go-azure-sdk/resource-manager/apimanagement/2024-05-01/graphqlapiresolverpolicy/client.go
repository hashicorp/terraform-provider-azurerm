package graphqlapiresolverpolicy

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GraphQLApiResolverPolicyClient struct {
	Client *resourcemanager.Client
}

func NewGraphQLApiResolverPolicyClientWithBaseURI(sdkApi sdkEnv.Api) (*GraphQLApiResolverPolicyClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "graphqlapiresolverpolicy", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating GraphQLApiResolverPolicyClient: %+v", err)
	}

	return &GraphQLApiResolverPolicyClient{
		Client: client,
	}, nil
}
