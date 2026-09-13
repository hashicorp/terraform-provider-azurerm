package prefixlistglobalrulestackresources

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrefixListGlobalRulestackResourcesClient struct {
	Client *resourcemanager.Client
}

func NewPrefixListGlobalRulestackResourcesClientWithBaseURI(sdkApi sdkEnv.Api) (*PrefixListGlobalRulestackResourcesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "prefixlistglobalrulestackresources", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PrefixListGlobalRulestackResourcesClient: %+v", err)
	}

	return &PrefixListGlobalRulestackResourcesClient{
		Client: client,
	}, nil
}
