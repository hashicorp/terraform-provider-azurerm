package localrulestacks

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LocalRulestacksClient struct {
	Client *resourcemanager.Client
}

func NewLocalRulestacksClientWithBaseURI(sdkApi sdkEnv.Api) (*LocalRulestacksClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "localrulestacks", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating LocalRulestacksClient: %+v", err)
	}

	return &LocalRulestacksClient{
		Client: client,
	}, nil
}
