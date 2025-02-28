package allpolicies

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AllPoliciesClient struct {
	Client *resourcemanager.Client
}

func NewAllPoliciesClientWithBaseURI(sdkApi sdkEnv.Api) (*AllPoliciesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "allpolicies", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AllPoliciesClient: %+v", err)
	}

	return &AllPoliciesClient{
		Client: client,
	}, nil
}
