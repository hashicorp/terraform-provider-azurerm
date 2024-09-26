package accesspolicies

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessPoliciesClient struct {
	Client *resourcemanager.Client
}

func NewAccessPoliciesClientWithBaseURI(sdkApi sdkEnv.Api) (*AccessPoliciesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "accesspolicies", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AccessPoliciesClient: %+v", err)
	}

	return &AccessPoliciesClient{
		Client: client,
	}, nil
}
