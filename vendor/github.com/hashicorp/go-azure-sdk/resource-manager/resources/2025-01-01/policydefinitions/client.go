package policydefinitions

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PolicyDefinitionsClient struct {
	Client *resourcemanager.Client
}

func NewPolicyDefinitionsClientWithBaseURI(sdkApi sdkEnv.Api) (*PolicyDefinitionsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "policydefinitions", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PolicyDefinitionsClient: %+v", err)
	}

	return &PolicyDefinitionsClient{
		Client: client,
	}, nil
}
