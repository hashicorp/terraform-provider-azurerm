// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package machinelearning_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2024-04-01/machinelearningcomputes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ComputeInstanceResource struct{}

func TestAccComputeInstance_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_compute_instance", "test")
	r := ComputeInstanceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccComputeInstance_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_compute_instance", "test")
	r := ComputeInstanceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ssh.0.username").Exists(),
				check.That(data.ResourceName).Key("ssh.0.port").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccComputeInstance_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_compute_instance", "test")
	r := ComputeInstanceResource{}

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

func TestAccComputeInstance_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_compute_instance", "test")
	r := ComputeInstanceResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.identitySystemAssignedUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.identityUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.identitySystemAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func (r ComputeInstanceResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	computeClient := client.MachineLearning.MachineLearningComputes
	id, err := machinelearningcomputes.ParseComputeID(state.ID)
	if err != nil {
		return nil, err
	}

	computeResource, err := computeClient.ComputeGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(computeResource.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Machine Learning Compute Cluster %q: %+v", state.ID, err)
	}
	return utils.Bool(computeResource.Model.Properties != nil), nil
}

func (r ComputeInstanceResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	var location string
	if !features.FourPointOhBeta() {
		location = "location = azurerm_resource_group.test.location"
	}
	return fmt.Sprintf(`
%s

resource "azurerm_machine_learning_compute_instance" "test" {
  name = "acctest%d"
  %s
  machine_learning_workspace_id = azurerm_machine_learning_workspace.test.id
  virtual_machine_size          = "STANDARD_DS2_V2"
  local_auth_enabled            = false
}
`, template, data.RandomIntOfLength(8), location)
}

func (r ComputeInstanceResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	var location string
	if !features.FourPointOhBeta() {
		location = "location = azurerm_resource_group.test.location"
	}
	return fmt.Sprintf(`
%s

variable "ssh_key" {
  default = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_security_group" "test" {
  name                = "test-nsg-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  security_rule {
    name                       = "test123"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "29876-44224"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }
}

resource "azurerm_subnet_network_security_group_association" "test" {
  subnet_id                 = azurerm_subnet.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}

resource "azurerm_private_dns_zone" "test" {
  name                = "privatelink.api.azureml.ms"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_zone_virtual_network_link" "test" {
  name                  = "test-vlink"
  resource_group_name   = azurerm_resource_group.test.name
  private_dns_zone_name = azurerm_private_dns_zone.test.name
  virtual_network_id    = azurerm_virtual_network.test.id
}

resource "azurerm_private_endpoint" "test" {
  name                = "test-pe-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.test.id
  private_service_connection {
    name                           = "test-mlworkspace-%d"
    private_connection_resource_id = azurerm_machine_learning_workspace.test.id
    subresource_names              = ["amlworkspace"]
    is_manual_connection           = false
  }
  private_dns_zone_group {
    name                 = "test"
    private_dns_zone_ids = [azurerm_private_dns_zone.test.id]
  }
}

resource "azurerm_machine_learning_compute_instance" "test" {
  name = "acctest%d"
  %s
  machine_learning_workspace_id = azurerm_machine_learning_workspace.test.id
  virtual_machine_size          = "STANDARD_DS2_V2"
  authorization_type            = "personal"
  node_public_ip_enabled        = false
  ssh {
    public_key = var.ssh_key
  }
  subnet_resource_id = azurerm_subnet.test.id
  description        = "this is desc"
  tags = {
    Label1 = "Value1"
  }
  depends_on = [
    azurerm_subnet_network_security_group_association.test,
    azurerm_private_endpoint.test
  ]
}
`, template, data.RandomIntOfLength(8), data.RandomIntOfLength(8), data.RandomIntOfLength(8), data.RandomIntOfLength(8), data.RandomIntOfLength(8), location)
}

func (r ComputeInstanceResource) requiresImport(data acceptance.TestData) string {
	var template string
	var location string

	if !features.FourPointOhBeta() {
		location = "location = azurerm_resource_group.test.location"
	}

	template = r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_machine_learning_compute_instance" "import" {
  name = azurerm_machine_learning_compute_instance.test.name
  %s
  machine_learning_workspace_id = azurerm_machine_learning_compute_instance.test.machine_learning_workspace_id
  virtual_machine_size          = "STANDARD_DS2_V2"
}
`, template, location)
}

func (r ComputeInstanceResource) identitySystemAssigned(data acceptance.TestData) string {
	template := r.template(data)
	var location string
	if !features.FourPointOhBeta() {
		location = "location = azurerm_resource_group.test.location"
	}
	return fmt.Sprintf(`
%s

resource "azurerm_machine_learning_compute_instance" "test" {
  name = "acctest%d"
  %s
  machine_learning_workspace_id = azurerm_machine_learning_workspace.test.id
  virtual_machine_size          = "STANDARD_DS2_V2"
  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomIntOfLength(8), location)
}

func (r ComputeInstanceResource) identityUserAssigned(data acceptance.TestData) string {
	template := r.template(data)
	var location string
	if !features.FourPointOhBeta() {
		location = "location = azurerm_resource_group.test.location"
	}
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_machine_learning_compute_instance" "test" {
  name = "acctest%d"
  %s
  machine_learning_workspace_id = azurerm_machine_learning_workspace.test.id
  virtual_machine_size          = "STANDARD_DS2_V2"
  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }
}
`, template, data.RandomInteger, data.RandomIntOfLength(8), location)
}

func (r ComputeInstanceResource) identitySystemAssignedUserAssigned(data acceptance.TestData) string {
	var location string
	template := r.template(data)

	if !features.FourPointOhBeta() {
		location = "location = azurerm_resource_group.test.location"
	}
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_machine_learning_compute_instance" "test" {
  name = "acctest%d"
  %s
  machine_learning_workspace_id = azurerm_machine_learning_workspace.test.id
  virtual_machine_size          = "STANDARD_DS2_V2"
  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }
}
`, template, data.RandomInteger, data.RandomIntOfLength(8), location)
}

func (r ComputeInstanceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ml-%[1]d"
  location = "%[2]s"
  tags = {
    "stage" = "test"
  }
}

resource "azurerm_application_insights" "test" {
  name                = "acctestai-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_key_vault" "test" {
  name                     = "acckv%[3]d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[4]d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_machine_learning_workspace" "test" {
  name                          = "acctest-MLW%[5]d"
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  application_insights_id       = azurerm_application_insights.test.id
  key_vault_id                  = azurerm_key_vault.test.id
  storage_account_id            = azurerm_storage_account.test.id
  public_network_access_enabled = true
  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger,
		data.RandomIntOfLength(15), data.RandomIntOfLength(16))
}
