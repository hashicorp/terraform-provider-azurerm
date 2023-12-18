package variable

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VariableClient struct {
	Client *resourcemanager.Client
}

func NewVariableClientWithBaseURI(sdkApi sdkEnv.Api) (*VariableClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "variable", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating VariableClient: %+v", err)
	}

	return &VariableClient{
		Client: client,
	}, nil
}
