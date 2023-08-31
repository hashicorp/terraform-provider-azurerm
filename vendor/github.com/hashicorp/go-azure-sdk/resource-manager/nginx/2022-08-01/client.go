package v2022_08_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2022-08-01/nginxcertificate"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2022-08-01/nginxconfiguration"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2022-08-01/nginxdeployment"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	NginxCertificate   *nginxcertificate.NginxCertificateClient
	NginxConfiguration *nginxconfiguration.NginxConfigurationClient
	NginxDeployment    *nginxdeployment.NginxDeploymentClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	nginxCertificateClient, err := nginxcertificate.NewNginxCertificateClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building NginxCertificate client: %+v", err)
	}
	configureFunc(nginxCertificateClient.Client)

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
		NginxCertificate:   nginxCertificateClient,
		NginxConfiguration: nginxConfigurationClient,
		NginxDeployment:    nginxDeploymentClient,
	}, nil
}
