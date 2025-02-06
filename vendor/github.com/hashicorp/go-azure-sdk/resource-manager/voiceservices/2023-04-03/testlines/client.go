package testlines

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TestLinesClient struct {
	Client *resourcemanager.Client
}

func NewTestLinesClientWithBaseURI(sdkApi sdkEnv.Api) (*TestLinesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "testlines", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating TestLinesClient: %+v", err)
	}

	return &TestLinesClient{
		Client: client,
	}, nil
}
