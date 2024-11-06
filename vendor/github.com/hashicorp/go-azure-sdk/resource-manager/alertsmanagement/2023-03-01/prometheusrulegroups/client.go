package prometheusrulegroups

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrometheusRuleGroupsClient struct {
	Client *resourcemanager.Client
}

func NewPrometheusRuleGroupsClientWithBaseURI(sdkApi sdkEnv.Api) (*PrometheusRuleGroupsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "prometheusrulegroups", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PrometheusRuleGroupsClient: %+v", err)
	}

	return &PrometheusRuleGroupsClient{
		Client: client,
	}, nil
}
