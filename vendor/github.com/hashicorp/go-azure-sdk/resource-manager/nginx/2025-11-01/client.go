package v2025_11_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2025-11-01/nginxcertificates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2025-11-01/nginxconfigurationresponses"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2025-11-01/nginxdeploymentapikeyresponses"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2025-11-01/nginxdeployments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2025-11-01/nginxdeploymentwafpolicies"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	NginxCertificates              *nginxcertificates.NginxCertificatesClient
	NginxConfigurationResponses    *nginxconfigurationresponses.NginxConfigurationResponsesClient
	NginxDeploymentApiKeyResponses *nginxdeploymentapikeyresponses.NginxDeploymentApiKeyResponsesClient
	NginxDeploymentWafPolicies     *nginxdeploymentwafpolicies.NginxDeploymentWafPoliciesClient
	NginxDeployments               *nginxdeployments.NginxDeploymentsClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	nginxCertificatesClient, err := nginxcertificates.NewNginxCertificatesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building NginxCertificates client: %+v", err)
	}
	configureFunc(nginxCertificatesClient.Client)

	nginxConfigurationResponsesClient, err := nginxconfigurationresponses.NewNginxConfigurationResponsesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building NginxConfigurationResponses client: %+v", err)
	}
	configureFunc(nginxConfigurationResponsesClient.Client)

	nginxDeploymentApiKeyResponsesClient, err := nginxdeploymentapikeyresponses.NewNginxDeploymentApiKeyResponsesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building NginxDeploymentApiKeyResponses client: %+v", err)
	}
	configureFunc(nginxDeploymentApiKeyResponsesClient.Client)

	nginxDeploymentWafPoliciesClient, err := nginxdeploymentwafpolicies.NewNginxDeploymentWafPoliciesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building NginxDeploymentWafPolicies client: %+v", err)
	}
	configureFunc(nginxDeploymentWafPoliciesClient.Client)

	nginxDeploymentsClient, err := nginxdeployments.NewNginxDeploymentsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building NginxDeployments client: %+v", err)
	}
	configureFunc(nginxDeploymentsClient.Client)

	return &Client{
		NginxCertificates:              nginxCertificatesClient,
		NginxConfigurationResponses:    nginxConfigurationResponsesClient,
		NginxDeploymentApiKeyResponses: nginxDeploymentApiKeyResponsesClient,
		NginxDeploymentWafPolicies:     nginxDeploymentWafPoliciesClient,
		NginxDeployments:               nginxDeploymentsClient,
	}, nil
}
