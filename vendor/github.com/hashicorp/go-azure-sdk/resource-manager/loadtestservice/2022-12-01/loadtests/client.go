package loadtests

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LoadTestsClient struct {
	Client *resourcemanager.Client
}

func NewLoadTestsClientWithBaseURI(sdkApi sdkEnv.Api) (*LoadTestsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "loadtests", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating LoadTestsClient: %+v", err)
	}

	return &LoadTestsClient{
		Client: client,
	}, nil
}
