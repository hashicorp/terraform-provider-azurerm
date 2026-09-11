package staticcidrs

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StaticCidrsClient struct {
	Client *resourcemanager.Client
}

func NewStaticCidrsClientWithBaseURI(sdkApi sdkEnv.Api) (*StaticCidrsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "staticcidrs", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating StaticCidrsClient: %+v", err)
	}

	return &StaticCidrsClient{
		Client: client,
	}, nil
}
