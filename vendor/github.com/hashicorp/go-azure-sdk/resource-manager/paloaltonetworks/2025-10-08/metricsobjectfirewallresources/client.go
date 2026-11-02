package metricsobjectfirewallresources

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MetricsObjectFirewallResourcesClient struct {
	Client *resourcemanager.Client
}

func NewMetricsObjectFirewallResourcesClientWithBaseURI(sdkApi sdkEnv.Api) (*MetricsObjectFirewallResourcesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "metricsobjectfirewallresources", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating MetricsObjectFirewallResourcesClient: %+v", err)
	}

	return &MetricsObjectFirewallResourcesClient{
		Client: client,
	}, nil
}
