// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package nginx_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-06-01-preview/nginxconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/nginx"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ConfigurationResource struct{}

func (a ConfigurationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := nginxconfiguration.ParseConfigurationID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Nginx.NginxConfiguration.ConfigurationsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func TestAccConfiguration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, nginx.ConfigurationResource{}.ResourceType(), "test")
	r := ConfigurationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("protected_file"),
	})
}

func TestAccConfiguration_withCertificate(t *testing.T) {
	data := acceptance.BuildTestData(t, nginx.ConfigurationResource{}.ResourceType(), "test")
	r := ConfigurationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withCertificate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccConfiguration_update(t *testing.T) {
	data := acceptance.BuildTestData(t, nginx.ConfigurationResource{}.ResourceType(), "test")
	r := ConfigurationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("protected_file"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("protected_file"),
	})
}

func TestAccConfiguration_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, nginx.ConfigurationResource{}.ResourceType(), "test")
	r := ConfigurationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (a ConfigurationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_nginx_configuration" "test" {
  nginx_deployment_id = azurerm_nginx_deployment.test.id
  root_file           = "/etc/nginx/nginx.conf"

  config_file {
    content      = local.config_content
    virtual_path = "/etc/nginx/nginx.conf"
  }

  protected_file {
    content      = local.protected_content
    virtual_path = "/opt/.htpasswd"
  }
}
`, a.template(data))
}

func (a ConfigurationResource) withCertificate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

locals {
  config_content = base64encode(<<-EOT
http {
    server {
      listen 443 ssl;
      ssl_certificate /etc/nginx/ssl/test.crt;
      ssl_certificate_key /etc/nginx/ssl/test.key;
      location / {
        return 200 "Hello World";
      }
    }
}
EOT
  )
}

resource "azurerm_nginx_certificate" "test" {
  name                     = "acctest"
  nginx_deployment_id      = azurerm_nginx_deployment.test.id
  key_virtual_path         = "/etc/nginx/ssl/test.key"
  certificate_virtual_path = "/etc/nginx/ssl/test.crt"
  key_vault_secret_id      = azurerm_key_vault_certificate.test.secret_id
}

resource "azurerm_nginx_configuration" "test" {
  nginx_deployment_id = azurerm_nginx_deployment.test.id
  root_file           = "/etc/nginx/nginx.conf"

  config_file {
    content      = local.config_content
    virtual_path = "/etc/nginx/nginx.conf"
  }

  depends_on = [azurerm_nginx_certificate.test]
}
`, CertificateResource{}.template(data))
}

func (a ConfigurationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_nginx_configuration" "import" {
  nginx_deployment_id = azurerm_nginx_configuration.test.nginx_deployment_id
  root_file           = azurerm_nginx_configuration.test.root_file
  config_file {
    content      = base64encode("http{}")
    virtual_path = "/"
  }
}
`, a.basic(data))
}

func (a ConfigurationResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_nginx_configuration" "test" {
  nginx_deployment_id = azurerm_nginx_deployment.test.id
  root_file           = "/etc/nginx/nginx.conf"

  config_file {
    content      = local.config_content
    virtual_path = "/etc/nginx/nginx.conf"
  }

  config_file {
    content      = local.sub_config_content
    virtual_path = "/etc/nginx/site/b.conf"
  }

  protected_file {
    content      = local.protected_content
    virtual_path = "/opt/.htpasswd"
  }
}
`, a.template(data))
}

func (a ConfigurationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%[1]d"
  location = "%[2]s"
}

locals {
  config_content = base64encode(<<-EOT
http {
    server {
        listen 80;
        location / {
            auth_basic "Protected Area";
            auth_basic_user_file /opt/.htpasswd;
            default_type text/html;
            return 200 '<!doctype html><html lang="en"><head></head><body>
                <div>this one will be updated</div>
                <div>at 10:38 am</div>
            </body></html>';
        }
        include site/*.conf;
    }
}
EOT
  )

  protected_content = base64encode(<<-EOT
user:$apr1$VeUA5kt.$IjjRk//8miRxDsZvD4daF1
EOT
  )

  sub_config_content = base64encode(<<-EOT
location /bbb {
	default_type text/html;
	return 200 '<!doctype html><html lang="en"><head></head><body>
		<div>this one will be updated</div>
		<div>at 10:38 am</div>
	</body></html>';
}
EOT
  )
}

resource "azurerm_public_ip" "test" {
  name                = "acctest%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  allocation_method   = "Static"
  sku                 = "Standard"

  tags = {
    environment = "Production"
  }
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "accsubnet%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
  delegation {
    name = "delegation"

    service_delegation {
      name = "NGINX.NGINXPLUS/nginxDeployments"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}

resource "azurerm_nginx_deployment" "test" {
  name                = "acctest-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "standardv2_Monthly"
  capacity            = 10
  location            = azurerm_resource_group.test.location

  diagnose_support_enabled = false

  frontend_public {
    ip_address = [azurerm_public_ip.test.id]
  }

  network_interface {
    subnet_id = azurerm_subnet.test.id
  }

  tags = {
    foo = "bar"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
