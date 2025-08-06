package metricalerts

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MetricAlertsClient struct {
	Client *resourcemanager.Client
}

func NewMetricAlertsClientWithBaseURI(sdkApi sdkEnv.Api) (*MetricAlertsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "metricalerts", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating MetricAlertsClient: %+v", err)
	}

	return &MetricAlertsClient{
		Client: client,
	}, nil
}
