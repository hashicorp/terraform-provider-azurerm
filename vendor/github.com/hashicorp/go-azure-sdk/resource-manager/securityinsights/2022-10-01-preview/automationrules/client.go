package automationrules

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutomationRulesClient struct {
	Client *resourcemanager.Client
}

func NewAutomationRulesClientWithBaseURI(sdkApi sdkEnv.Api) (*AutomationRulesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "automationrules", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AutomationRulesClient: %+v", err)
	}

	return &AutomationRulesClient{
		Client: client,
	}, nil
}
