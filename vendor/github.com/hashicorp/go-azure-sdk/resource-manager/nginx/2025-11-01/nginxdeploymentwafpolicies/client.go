package nginxdeploymentwafpolicies

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NginxDeploymentWafPoliciesClient struct {
	Client *resourcemanager.Client
}

func NewNginxDeploymentWafPoliciesClientWithBaseURI(sdkApi sdkEnv.Api) (*NginxDeploymentWafPoliciesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "nginxdeploymentwafpolicies", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NginxDeploymentWafPoliciesClient: %+v", err)
	}

	return &NginxDeploymentWafPoliciesClient{
		Client: client,
	}, nil
}
