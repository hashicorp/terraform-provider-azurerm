package authenticationpolicies

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthenticationPoliciesClient struct {
	Client *resourcemanager.Client
}

func NewAuthenticationPoliciesClientWithBaseURI(sdkApi sdkEnv.Api) (*AuthenticationPoliciesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "authenticationpolicies", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AuthenticationPoliciesClient: %+v", err)
	}

	return &AuthenticationPoliciesClient{
		Client: client,
	}, nil
}
