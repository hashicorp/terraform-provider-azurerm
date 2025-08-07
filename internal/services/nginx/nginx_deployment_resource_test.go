// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package nginx_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-11-01-preview/nginxdeployment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/nginx"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DeploymentResource struct{}

func (a DeploymentResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := nginxdeployment.ParseNginxDeploymentID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Nginx.NginxDeployment.DeploymentsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving Deployment %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func TestAccNginxDeployment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, nginx.DeploymentResource{}.ResourceType(), "test")
	r := DeploymentResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("capacity").HasValue("10"),
				check.That(data.ResourceName).Key("email").HasValue("test@test.com"),
				check.That(data.ResourceName).Key("dataplane_api_endpoint").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNginxDeployment_update(t *testing.T) {
	data := acceptance.BuildTestData(t, nginx.DeploymentResource{}.ResourceType(), "test")
	r := DeploymentResource{}
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

func TestAccNginxDeployment_updateNetworkInterface(t *testing.T) {
	data := acceptance.BuildTestData(t, nginx.DeploymentResource{}.ResourceType(), "test")
	r := DeploymentResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicPrivate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateNetworkInterface(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_interface.0.subnet_id").MatchesRegex(regexp.MustCompile(`^/subscriptions/[\w-]+/resourceGroups/[\w-]+/providers/Microsoft.Network/virtualNetworks/[\w-]+/subnets/subnet2-[\w-]+$`)),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNginxDeployment_updateFrontendPrivate(t *testing.T) {
	data := acceptance.BuildTestData(t, nginx.DeploymentResource{}.ResourceType(), "test")
	r := DeploymentResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicPrivate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateFrontendPrivate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("frontend_private.0.ip_address").HasValue("10.0.2.11"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNginxDeployment_updateFrontendPublic(t *testing.T) {
	data := acceptance.BuildTestData(t, nginx.DeploymentResource{}.ResourceType(), "test")
	r := DeploymentResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateFrontendPublic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("frontend_public.0.ip_address.0").MatchesRegex(regexp.MustCompile(`^/subscriptions/[\w-]+/resourceGroups/[\w-]+/providers/Microsoft.Network/publicIPAddresses/acctest2-[\w-]+$`)),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNginxDeployment_systemAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, nginx.DeploymentResource{}.ResourceType(), "test")
	r := DeploymentResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.systemAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNginxDeployment_userAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, nginx.DeploymentResource{}.ResourceType(), "test")
	r := DeploymentResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.userAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (a DeploymentResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_nginx_deployment" "test" {
  name                      = "acctest-%[2]d"
  resource_group_name       = azurerm_resource_group.test.name
  sku                       = "standardv2_Monthly"
  location                  = azurerm_resource_group.test.location
  diagnose_support_enabled  = false
  automatic_upgrade_channel = "stable"

  frontend_public {
    ip_address = [azurerm_public_ip.test.id]
  }

  network_interface {
    subnet_id = azurerm_subnet.test.id
  }

  capacity = 10

  email = "test@test.com"

  tags = {
    foo = "bar"
  }
}
`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a DeploymentResource) basicPrivate(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_nginx_deployment" "test" {
  name                      = "acctest-%[2]d"
  resource_group_name       = azurerm_resource_group.test.name
  sku                       = "standardv2_Monthly"
  location                  = azurerm_resource_group.test.location
  diagnose_support_enabled  = false
  automatic_upgrade_channel = "stable"

  frontend_private {
    allocation_method = "Static"
    ip_address        = "10.0.2.10"
    subnet_id         = azurerm_subnet.test.id
  }

  network_interface {
    subnet_id = azurerm_subnet.test.id
  }

  capacity = 10

  email = "test@test.com"

  tags = {
    foo = "bar"
  }
}
`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a DeploymentResource) basicAutoscaling(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_nginx_deployment" "test" {
  name                      = "acctest-%[2]d"
  resource_group_name       = azurerm_resource_group.test.name
  sku                       = "standardv2_Monthly"
  location                  = azurerm_resource_group.test.location
  diagnose_support_enabled  = false
  automatic_upgrade_channel = "stable"

  frontend_public {
    ip_address = [azurerm_public_ip.test.id]
  }

  network_interface {
    subnet_id = azurerm_subnet.test.id
  }

  auto_scale_profile {
    name         = "test"
    min_capacity = 10
    max_capacity = 30
  }

  email = "test@test.com"

  tags = {
    foo = "bar"
  }

  lifecycle {
    ignore_changes = [
      capacity,
    ]
  }
}
`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a DeploymentResource) basicAutoscaling_update(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_nginx_deployment" "test" {
  name                      = "acctest-%[2]d"
  resource_group_name       = azurerm_resource_group.test.name
  sku                       = "standardv2_Monthly"
  location                  = azurerm_resource_group.test.location
  diagnose_support_enabled  = false
  automatic_upgrade_channel = "stable"

  frontend_public {
    ip_address = [azurerm_public_ip.test.id]
  }

  network_interface {
    subnet_id = azurerm_subnet.test.id
  }

  auto_scale_profile {
    name         = "test"
    min_capacity = 10
    max_capacity = 20
  }

  email = "test@test.com"

  tags = {
    foo = "bar"
  }

  lifecycle {
    ignore_changes = [
      capacity,
    ]
  }
}
`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a DeploymentResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_nginx_deployment" "test" {
  name                     = "acctest-%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  sku                      = "standardv2_Monthly"
  location                 = azurerm_resource_group.test.location
  diagnose_support_enabled = false

  frontend_public {
    ip_address = [azurerm_public_ip.test.id]
  }

  network_interface {
    subnet_id = azurerm_subnet.test.id
  }

  capacity = 20

  email = "testing@test.com"

  tags = {
    foo = "bar2"
  }
}
`, a.template(data), data.RandomInteger)
}

func (a DeploymentResource) updateNetworkInterface(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_nginx_deployment" "test" {
  name                      = "acctest-%[2]d"
  resource_group_name       = azurerm_resource_group.test.name
  sku                       = "standardv2_Monthly"
  location                  = azurerm_resource_group.test.location
  diagnose_support_enabled  = false
  automatic_upgrade_channel = "stable"

  frontend_private {
    allocation_method = "Static"
    ip_address        = "10.0.2.10"
    subnet_id         = azurerm_subnet.test.id
  }

  network_interface {
    subnet_id = azurerm_subnet.test2.id
  }

  capacity = 10

  email = "testing@test.com"

  tags = {
    foo = "bar2"
  }
}
`, a.template(data), data.RandomInteger)
}

func (a DeploymentResource) updateFrontendPrivate(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_nginx_deployment" "test" {
  name                      = "acctest-%[2]d"
  resource_group_name       = azurerm_resource_group.test.name
  sku                       = "standardv2_Monthly"
  location                  = azurerm_resource_group.test.location
  diagnose_support_enabled  = false
  automatic_upgrade_channel = "stable"

  frontend_private {
    allocation_method = "Static"
    ip_address        = "10.0.2.11"
    subnet_id         = azurerm_subnet.test.id
  }

  network_interface {
    subnet_id = azurerm_subnet.test.id
  }

  capacity = 10

  email = "test@test.com"

  tags = {
    foo = "bar"
  }
}
`, a.template(data), data.RandomInteger)
}

func (a DeploymentResource) updateFrontendPublic(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_nginx_deployment" "test" {
  name                      = "acctest-%[2]d"
  resource_group_name       = azurerm_resource_group.test.name
  sku                       = "standardv2_Monthly"
  location                  = azurerm_resource_group.test.location
  diagnose_support_enabled  = false
  automatic_upgrade_channel = "stable"

  frontend_public {
    ip_address = [azurerm_public_ip.test2.id]
  }

  network_interface {
    subnet_id = azurerm_subnet.test.id
  }

  capacity = 10

  email = "testing@test.com"

  tags = {
    foo = "bar2"
  }
}
`, a.template(data), data.RandomInteger)
}

func (a DeploymentResource) systemAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_nginx_deployment" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "standardv2_Monthly"
  location            = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }

  frontend_public {
    ip_address = [azurerm_public_ip.test.id]
  }

  network_interface {
    subnet_id = azurerm_subnet.test.id
  }

  capacity = 10

  email = "test@test.com"
}
`, a.template(data), data.RandomInteger)
}

func (a DeploymentResource) userAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acct-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_nginx_deployment" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "standardv2_Monthly"
  location            = azurerm_resource_group.test.location

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  frontend_public {
    ip_address = [azurerm_public_ip.test.id]
  }

  network_interface {
    subnet_id = azurerm_subnet.test.id
  }

  capacity = 10

  email = "test@test.com"
}
`, a.template(data), data.RandomInteger)
}

func (a DeploymentResource) basicNginxAppProtect(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_nginx_deployment" "test" {
  name                      = "acctest-%[2]d"
  resource_group_name       = azurerm_resource_group.test.name
  sku                       = "standardv2_Monthly"
  location                  = azurerm_resource_group.test.location
  diagnose_support_enabled  = false
  automatic_upgrade_channel = "stable"

  frontend_public {
    ip_address = [azurerm_public_ip.test.id]
  }

  network_interface {
    subnet_id = azurerm_subnet.test.id
  }

  web_application_firewall {
    activation_state_enabled = true
  }

  capacity = 20

  email = "test@test.com"
}
`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a DeploymentResource) template(data acceptance.TestData) string {
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

resource "azurerm_public_ip" "test2" {
  name                = "acctest2-%[1]d"
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
  name                 = "subnet-%[1]d"
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

resource "azurerm_subnet" "test2" {
  name                 = "subnet2-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.3.0/24"]
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
`, data.RandomInteger, data.Locations.Primary)
}

func TestAccNginxDeployment_autoscaling(t *testing.T) {
	data := acceptance.BuildTestData(t, nginx.DeploymentResource{}.ResourceType(), "test")
	r := DeploymentResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicAutoscaling(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auto_scale_profile.0.name").HasValue("test"),
				check.That(data.ResourceName).Key("auto_scale_profile.0.min_capacity").HasValue("10"),
				check.That(data.ResourceName).Key("auto_scale_profile.0.max_capacity").HasValue("30"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicAutoscaling_update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auto_scale_profile.0.name").HasValue("test"),
				check.That(data.ResourceName).Key("auto_scale_profile.0.min_capacity").HasValue("10"),
				check.That(data.ResourceName).Key("auto_scale_profile.0.max_capacity").HasValue("20"),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auto_scale_profile.0.name").DoesNotExist(),
				check.That(data.ResourceName).Key("auto_scale_profile.0.min_capacity").DoesNotExist(),
				check.That(data.ResourceName).Key("auto_scale_profile.0.max_capacity").DoesNotExist(),
			),
		},
		data.ImportStep(),
	})
}
