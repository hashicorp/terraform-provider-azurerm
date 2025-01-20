package ddoscustompolicies

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DdosCustomPoliciesClient struct {
	Client *resourcemanager.Client
}

func NewDdosCustomPoliciesClientWithBaseURI(sdkApi sdkEnv.Api) (*DdosCustomPoliciesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "ddoscustompolicies", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DdosCustomPoliciesClient: %+v", err)
	}

	return &DdosCustomPoliciesClient{
		Client: client,
	}, nil
}
