package fqdnlistglobalrulestackresources

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FqdnListGlobalRulestackResourcesClient struct {
	Client *resourcemanager.Client
}

func NewFqdnListGlobalRulestackResourcesClientWithBaseURI(sdkApi sdkEnv.Api) (*FqdnListGlobalRulestackResourcesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "fqdnlistglobalrulestackresources", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating FqdnListGlobalRulestackResourcesClient: %+v", err)
	}

	return &FqdnListGlobalRulestackResourcesClient{
		Client: client,
	}, nil
}
