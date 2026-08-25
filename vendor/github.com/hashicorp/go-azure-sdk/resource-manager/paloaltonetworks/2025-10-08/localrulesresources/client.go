package localrulesresources

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LocalRulesResourcesClient struct {
	Client *resourcemanager.Client
}

func NewLocalRulesResourcesClientWithBaseURI(sdkApi sdkEnv.Api) (*LocalRulesResourcesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "localrulesresources", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating LocalRulesResourcesClient: %+v", err)
	}

	return &LocalRulesResourcesClient{
		Client: client,
	}, nil
}
