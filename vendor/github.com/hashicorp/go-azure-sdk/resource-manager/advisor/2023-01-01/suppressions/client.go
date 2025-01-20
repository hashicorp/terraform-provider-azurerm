package suppressions

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SuppressionsClient struct {
	Client *resourcemanager.Client
}

func NewSuppressionsClientWithBaseURI(sdkApi sdkEnv.Api) (*SuppressionsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "suppressions", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SuppressionsClient: %+v", err)
	}

	return &SuppressionsClient{
		Client: client,
	}, nil
}
