package policyrestriction

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PolicyRestrictionClient struct {
	Client *resourcemanager.Client
}

func NewPolicyRestrictionClientWithBaseURI(sdkApi sdkEnv.Api) (*PolicyRestrictionClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "policyrestriction", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PolicyRestrictionClient: %+v", err)
	}

	return &PolicyRestrictionClient{
		Client: client,
	}, nil
}
