package alertruletemplates

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertRuleTemplatesClient struct {
	Client *resourcemanager.Client
}

func NewAlertRuleTemplatesClientWithBaseURI(sdkApi sdkEnv.Api) (*AlertRuleTemplatesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "alertruletemplates", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AlertRuleTemplatesClient: %+v", err)
	}

	return &AlertRuleTemplatesClient{
		Client: client,
	}, nil
}
