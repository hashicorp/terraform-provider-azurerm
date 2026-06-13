package nginxconfigurationresponses

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NginxConfigurationResponsesClient struct {
	Client *resourcemanager.Client
}

func NewNginxConfigurationResponsesClientWithBaseURI(sdkApi sdkEnv.Api) (*NginxConfigurationResponsesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "nginxconfigurationresponses", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NginxConfigurationResponsesClient: %+v", err)
	}

	return &NginxConfigurationResponsesClient{
		Client: client,
	}, nil
}
