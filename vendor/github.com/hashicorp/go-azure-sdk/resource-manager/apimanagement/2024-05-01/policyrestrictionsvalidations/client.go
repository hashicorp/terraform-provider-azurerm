package policyrestrictionsvalidations

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PolicyRestrictionsValidationsClient struct {
	Client *resourcemanager.Client
}

func NewPolicyRestrictionsValidationsClientWithBaseURI(sdkApi sdkEnv.Api) (*PolicyRestrictionsValidationsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "policyrestrictionsvalidations", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PolicyRestrictionsValidationsClient: %+v", err)
	}

	return &PolicyRestrictionsValidationsClient{
		Client: client,
	}, nil
}
