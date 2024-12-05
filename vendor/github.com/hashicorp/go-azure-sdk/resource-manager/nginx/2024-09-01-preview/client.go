package v2024_09_01_preview

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-09-01-preview/nginxapikey"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-09-01-preview/nginxcertificate"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-09-01-preview/nginxconfiguration"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-09-01-preview/nginxconfigurationanalysis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-09-01-preview/nginxdeployment"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	NginxApiKey                *nginxapikey.NginxApiKeyClient
	NginxCertificate           *nginxcertificate.NginxCertificateClient
	NginxConfiguration         *nginxconfiguration.NginxConfigurationClient
	NginxConfigurationAnalysis *nginxconfigurationanalysis.NginxConfigurationAnalysisClient
	NginxDeployment            *nginxdeployment.NginxDeploymentClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	nginxApiKeyClient, err := nginxapikey.NewNginxApiKeyClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building NginxApiKey client: %+v", err)
	}
	configureFunc(nginxApiKeyClient.Client)

	nginxCertificateClient, err := nginxcertificate.NewNginxCertificateClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building NginxCertificate client: %+v", err)
	}
	configureFunc(nginxCertificateClient.Client)

	nginxConfigurationAnalysisClient, err := nginxconfigurationanalysis.NewNginxConfigurationAnalysisClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building NginxConfigurationAnalysis client: %+v", err)
	}
	configureFunc(nginxConfigurationAnalysisClient.Client)

	nginxConfigurationClient, err := nginxconfiguration.NewNginxConfigurationClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building NginxConfiguration client: %+v", err)
	}
	configureFunc(nginxConfigurationClient.Client)

	nginxDeploymentClient, err := nginxdeployment.NewNginxDeploymentClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building NginxDeployment client: %+v", err)
	}
	configureFunc(nginxDeploymentClient.Client)

	return &Client{
		NginxApiKey:                nginxApiKeyClient,
		NginxCertificate:           nginxCertificateClient,
		NginxConfiguration:         nginxConfigurationClient,
		NginxConfigurationAnalysis: nginxConfigurationAnalysisClient,
		NginxDeployment:            nginxDeploymentClient,
	}, nil
}
