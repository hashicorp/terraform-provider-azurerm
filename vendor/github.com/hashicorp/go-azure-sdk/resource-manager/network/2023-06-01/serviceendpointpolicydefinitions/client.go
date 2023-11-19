package serviceendpointpolicydefinitions

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceEndpointPolicyDefinitionsClient struct {
	Client *resourcemanager.Client
}

func NewServiceEndpointPolicyDefinitionsClientWithBaseURI(sdkApi sdkEnv.Api) (*ServiceEndpointPolicyDefinitionsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "serviceendpointpolicydefinitions", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ServiceEndpointPolicyDefinitionsClient: %+v", err)
	}

	return &ServiceEndpointPolicyDefinitionsClient{
		Client: client,
	}, nil
}
