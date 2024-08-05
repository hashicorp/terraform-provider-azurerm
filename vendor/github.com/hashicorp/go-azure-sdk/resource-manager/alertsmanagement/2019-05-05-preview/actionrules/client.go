package actionrules

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActionRulesClient struct {
	Client *resourcemanager.Client
}

func NewActionRulesClientWithBaseURI(sdkApi sdkEnv.Api) (*ActionRulesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "actionrules", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ActionRulesClient: %+v", err)
	}

	return &ActionRulesClient{
		Client: client,
	}, nil
}
