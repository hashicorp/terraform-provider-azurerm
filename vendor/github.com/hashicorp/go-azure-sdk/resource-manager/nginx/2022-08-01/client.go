package v2022_08_01

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2022-08-01/nginxcertificate"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2022-08-01/nginxconfiguration"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2022-08-01/nginxdeployment"
)

type Client struct {
	NginxCertificate   *nginxcertificate.NginxCertificateClient
	NginxConfiguration *nginxconfiguration.NginxConfigurationClient
	NginxDeployment    *nginxdeployment.NginxDeploymentClient
}

func NewClientWithBaseURI(endpoint string, configureAuthFunc func(c *autorest.Client)) Client {

	nginxCertificateClient := nginxcertificate.NewNginxCertificateClientWithBaseURI(endpoint)
	configureAuthFunc(&nginxCertificateClient.Client)

	nginxConfigurationClient := nginxconfiguration.NewNginxConfigurationClientWithBaseURI(endpoint)
	configureAuthFunc(&nginxConfigurationClient.Client)

	nginxDeploymentClient := nginxdeployment.NewNginxDeploymentClientWithBaseURI(endpoint)
	configureAuthFunc(&nginxDeploymentClient.Client)

	return Client{
		NginxCertificate:   &nginxCertificateClient,
		NginxConfiguration: &nginxConfigurationClient,
		NginxDeployment:    &nginxDeploymentClient,
	}
}
