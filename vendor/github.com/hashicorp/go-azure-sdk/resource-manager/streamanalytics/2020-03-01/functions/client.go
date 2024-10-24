package functions

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FunctionsClient struct {
	Client *resourcemanager.Client
}

func NewFunctionsClientWithBaseURI(sdkApi sdkEnv.Api) (*FunctionsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "functions", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating FunctionsClient: %+v", err)
	}

	return &FunctionsClient{
		Client: client,
	}, nil
}
