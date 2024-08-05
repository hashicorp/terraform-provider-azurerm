package globalrulestack

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GlobalRulestackClient struct {
	Client *resourcemanager.Client
}

func NewGlobalRulestackClientWithBaseURI(sdkApi sdkEnv.Api) (*GlobalRulestackClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "globalrulestack", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating GlobalRulestackClient: %+v", err)
	}

	return &GlobalRulestackClient{
		Client: client,
	}, nil
}
