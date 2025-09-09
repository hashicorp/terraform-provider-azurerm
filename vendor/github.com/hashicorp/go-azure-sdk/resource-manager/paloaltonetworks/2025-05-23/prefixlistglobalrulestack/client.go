package prefixlistglobalrulestack

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrefixListGlobalRulestackClient struct {
	Client *resourcemanager.Client
}

func NewPrefixListGlobalRulestackClientWithBaseURI(sdkApi sdkEnv.Api) (*PrefixListGlobalRulestackClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "prefixlistglobalrulestack", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PrefixListGlobalRulestackClient: %+v", err)
	}

	return &PrefixListGlobalRulestackClient{
		Client: client,
	}, nil
}
