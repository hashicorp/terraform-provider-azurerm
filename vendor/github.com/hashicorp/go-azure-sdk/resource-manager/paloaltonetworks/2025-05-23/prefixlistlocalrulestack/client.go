package prefixlistlocalrulestack

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrefixListLocalRulestackClient struct {
	Client *resourcemanager.Client
}

func NewPrefixListLocalRulestackClientWithBaseURI(sdkApi sdkEnv.Api) (*PrefixListLocalRulestackClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "prefixlistlocalrulestack", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PrefixListLocalRulestackClient: %+v", err)
	}

	return &PrefixListLocalRulestackClient{
		Client: client,
	}, nil
}
