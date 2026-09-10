package prerulesresources

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PreRulesResourcesClient struct {
	Client *resourcemanager.Client
}

func NewPreRulesResourcesClientWithBaseURI(sdkApi sdkEnv.Api) (*PreRulesResourcesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "prerulesresources", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PreRulesResourcesClient: %+v", err)
	}

	return &PreRulesResourcesClient{
		Client: client,
	}, nil
}
