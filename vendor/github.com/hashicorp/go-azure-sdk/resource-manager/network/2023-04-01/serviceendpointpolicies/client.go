package serviceendpointpolicies

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceEndpointPoliciesClient struct {
	Client *resourcemanager.Client
}

func NewServiceEndpointPoliciesClientWithBaseURI(sdkApi sdkEnv.Api) (*ServiceEndpointPoliciesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "serviceendpointpolicies", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ServiceEndpointPoliciesClient: %+v", err)
	}

	return &ServiceEndpointPoliciesClient{
		Client: client,
	}, nil
}
