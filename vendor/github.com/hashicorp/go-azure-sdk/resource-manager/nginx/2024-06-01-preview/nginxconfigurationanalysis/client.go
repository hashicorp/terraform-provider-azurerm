package nginxconfigurationanalysis

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NginxConfigurationAnalysisClient struct {
	Client *resourcemanager.Client
}

func NewNginxConfigurationAnalysisClientWithBaseURI(sdkApi sdkEnv.Api) (*NginxConfigurationAnalysisClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "nginxconfigurationanalysis", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NginxConfigurationAnalysisClient: %+v", err)
	}

	return &NginxConfigurationAnalysisClient{
		Client: client,
	}, nil
}
