package automationaccount

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutomationAccountClient struct {
	Client *resourcemanager.Client
}

func NewAutomationAccountClientWithBaseURI(sdkApi sdkEnv.Api) (*AutomationAccountClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "automationaccount", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AutomationAccountClient: %+v", err)
	}

	return &AutomationAccountClient{
		Client: client,
	}, nil
}
