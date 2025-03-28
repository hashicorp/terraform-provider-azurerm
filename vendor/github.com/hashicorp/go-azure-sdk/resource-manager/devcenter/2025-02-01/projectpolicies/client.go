package projectpolicies

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProjectPoliciesClient struct {
	Client *resourcemanager.Client
}

func NewProjectPoliciesClientWithBaseURI(sdkApi sdkEnv.Api) (*ProjectPoliciesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "projectpolicies", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ProjectPoliciesClient: %+v", err)
	}

	return &ProjectPoliciesClient{
		Client: client,
	}, nil
}
