package hdinsights

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HdinsightsClient struct {
	Client *resourcemanager.Client
}

func NewHdinsightsClientWithBaseURI(sdkApi sdkEnv.Api) (*HdinsightsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "hdinsights", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating HdinsightsClient: %+v", err)
	}

	return &HdinsightsClient{
		Client: client,
	}, nil
}
