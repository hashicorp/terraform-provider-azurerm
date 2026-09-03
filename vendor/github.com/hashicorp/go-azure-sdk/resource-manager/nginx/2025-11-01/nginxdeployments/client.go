package nginxdeployments

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NginxDeploymentsClient struct {
	Client *resourcemanager.Client
}

func NewNginxDeploymentsClientWithBaseURI(sdkApi sdkEnv.Api) (*NginxDeploymentsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "nginxdeployments", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NginxDeploymentsClient: %+v", err)
	}

	return &NginxDeploymentsClient{
		Client: client,
	}, nil
}
