package nginx_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/nginx"

	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2022-08-01/nginxconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ConfigurationResource struct{}

func (a ConfigurationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := nginxconfiguration.ParseConfigurationID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Nginx.NginxConfiguration.ConfigurationsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving Configuration %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
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
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (a ConfigurationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_nginx_configuration" "test" {
  nginx_deployment_id = azurerm_nginx_deployment.test.id
  root_file           = "/etc/nginx/nginx.conf"

  config_file {
    content      = "aHR0cCB7DQogICAgc2VydmVyIHsNCiAgICAgICAgbGlzdGVuIDgwOw0KICAgICAgICBsb2NhdGlvbiAvIHsNCiAgICAgICAgICAgIGRlZmF1bHRfdHlwZSB0ZXh0L2h0bWw7DQogICAgICAgICAgICByZXR1cm4gMjAwICc8IWRvY3R5cGUgaHRtbD48aHRtbCBsYW5nPSJlbiI+PGhlYWQ+PC9oZWFkPjxib2R5Pg0KICAgICAgICAgICAgICAgIDxkaXY+dGhpcyBvbmUgd2lsbCBiZSB1cGRhdGVkPC9kaXY+DQogICAgICAgICAgICAgICAgPGRpdj5hdCAxMDozOCBhbTwvZGl2Pg0KICAgICAgICAgICAgPC9ib2R5PjwvaHRtbD4nOw0KICAgICAgICB9DQogICAgICAgIGluY2x1ZGUgc2l0ZS8qLmNvbmY7DQogICAgfQ0KfQ=="
    virtual_path = "/etc/nginx/nginx.conf"
  }
}
`, a.template(data))
}

func (a ConfigurationResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_nginx_configuration" "test" {
  nginx_deployment_id = azurerm_nginx_deployment.test.id
  root_file           = "/etc/nginx/nginx.conf"

  config_file {
    content      = "aHR0cCB7DQogICAgc2VydmVyIHsNCiAgICAgICAgbGlzdGVuIDgwOw0KICAgICAgICBsb2NhdGlvbiAvIHsNCiAgICAgICAgICAgIGRlZmF1bHRfdHlwZSB0ZXh0L2h0bWw7DQogICAgICAgICAgICByZXR1cm4gMjAwICc8IWRvY3R5cGUgaHRtbD48aHRtbCBsYW5nPSJlbiI+PGhlYWQ+PC9oZWFkPjxib2R5Pg0KICAgICAgICAgICAgICAgIDxkaXY+dGhpcyBvbmUgd2lsbCBiZSB1cGRhdGVkPC9kaXY+DQogICAgICAgICAgICAgICAgPGRpdj5hdCAxMDozOCBhbTwvZGl2Pg0KICAgICAgICAgICAgPC9ib2R5PjwvaHRtbD4nOw0KICAgICAgICB9DQogICAgICAgIGluY2x1ZGUgc2l0ZS8qLmNvbmY7DQogICAgfQ0KfQ=="
    virtual_path = "/etc/nginx/nginx.conf"
  }

  config_file {
    content      = "DQogICAgICAgIGxvY2F0aW9uIC9iYmIgew0KICAgICAgICAgICAgZGVmYXVsdF90eXBlIHRleHQvaHRtbDsNCiAgICAgICAgICAgIHJldHVybiAyMDAgJzwhZG9jdHlwZSBodG1sPjxodG1sIGxhbmc9ImVuIj48aGVhZD48L2hlYWQ+PGJvZHk+DQogICAgICAgICAgICAgICAgPGRpdj50aGlzIG9uZSB3aWxsIGJlIHVwZGF0ZWQ8L2Rpdj4NCiAgICAgICAgICAgICAgICA8ZGl2PmF0IDEwOjM4IGFtPC9kaXY+DQogICAgICAgICAgICA8L2JvZHk+PC9odG1sPic7DQogICAgICAgIH0NCg=="
    virtual_path = "/etc/nginx/site/b.conf"
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
  sku                 = "publicpreview_Monthly_gmz7xq9ge3py"
  location            = azurerm_resource_group.test.location

  //message: "Conflict managed resource group name: tenant: -91a, subscription xxx, resource group example."
  managed_resource_group   = "accmr%[1]d"
  diagnose_support_enabled = true

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
