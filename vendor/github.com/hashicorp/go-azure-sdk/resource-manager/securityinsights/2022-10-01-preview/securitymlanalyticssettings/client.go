package securitymlanalyticssettings

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityMLAnalyticsSettingsClient struct {
	Client *resourcemanager.Client
}

func NewSecurityMLAnalyticsSettingsClientWithBaseURI(sdkApi sdkEnv.Api) (*SecurityMLAnalyticsSettingsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "securitymlanalyticssettings", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SecurityMLAnalyticsSettingsClient: %+v", err)
	}

	return &SecurityMLAnalyticsSettingsClient{
		Client: client,
	}, nil
}
