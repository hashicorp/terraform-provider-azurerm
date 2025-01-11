package skus

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SKUsClient struct {
	Client *resourcemanager.Client
}

func NewSKUsClientWithBaseURI(sdkApi sdkEnv.Api) (*SKUsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "skus", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SKUsClient: %+v", err)
	}

	return &SKUsClient{
		Client: client,
	}, nil
}
