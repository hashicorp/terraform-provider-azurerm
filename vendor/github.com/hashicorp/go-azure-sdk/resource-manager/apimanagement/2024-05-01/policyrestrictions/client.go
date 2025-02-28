package policyrestrictions

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PolicyRestrictionsClient struct {
	Client *resourcemanager.Client
}

func NewPolicyRestrictionsClientWithBaseURI(sdkApi sdkEnv.Api) (*PolicyRestrictionsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "policyrestrictions", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PolicyRestrictionsClient: %+v", err)
	}

	return &PolicyRestrictionsClient{
		Client: client,
	}, nil
}
