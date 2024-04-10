package nginxdeployment

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NginxDeploymentClient struct {
	Client *resourcemanager.Client
}

func NewNginxDeploymentClientWithBaseURI(sdkApi sdkEnv.Api) (*NginxDeploymentClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "nginxdeployment", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NginxDeploymentClient: %+v", err)
	}

	return &NginxDeploymentClient{
		Client: client,
	}, nil
}
