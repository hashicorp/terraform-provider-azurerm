// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicefabricmanaged_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicefabricmanagedcluster/2021-05-01/managedcluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ClusterResource struct{}

func TestAccServiceFabricManagedCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_managed_cluster", "test")
	r := ClusterResource{}
	nodeTypeData1 := r.nodeType("test1", true, 130)
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, nodeTypeData1),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.Test").HasValue("value")),
		},
		data.ImportStep("password"),
	})
}

func TestAccServiceFabricManagedCluster_withCustomSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_managed_cluster", "test")
	r := ClusterResource{}
	nodeTypeData1 := r.nodeType("test1", true, 130)
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withCustomSettings(data, nodeTypeData1),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.Test").HasValue("value")),
		},
		data.ImportStep("password"),
	})
}

func TestAccServiceFabricManagedCluster_importError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_managed_cluster", "test")

	r := ClusterResource{}
	nodeTypeData1 := r.nodeType("test1", true, 130)
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, nodeTypeData1),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.Test").HasValue("value")),
		},
		{
			Config:      r.requiresImport(data, nodeTypeData1),
			ExpectError: acceptance.RequiresImportError(data.ResourceType),
		},
	})
}

func TestAccServiceFabricManagedCluster_full(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_managed_cluster", "test")
	r := ClusterResource{}
	nodeTypeData1 := r.nodeType("test1", true, 130)
	nodeTypeData1Altered := r.nodeType("test1", true, 140)
	nodeTypeData2 := r.nodeType("test2", false, 130)
	nodeTypeDataBoth := fmt.Sprintf("%s\n%s", nodeTypeData1, nodeTypeData2)
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, nodeTypeData1),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.Test").HasValue("value")),
		},
		data.ImportStep("password"),
		{
			Config: r.basic(data, nodeTypeDataBoth),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("node_type.#").HasValue("2"),
				check.That(data.ResourceName).Key("node_type.0.name").HasValue("test1"),
				check.That(data.ResourceName).Key("node_type.1.name").HasValue("test2"),
			),
		},
		{
			Config: r.basic(data, nodeTypeData1),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("node_type.#").HasValue("1"),
				check.That(data.ResourceName).Key("node_type.0.name").HasValue("test1")),
		},
		{
			Config: r.basic(data, nodeTypeData1Altered),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("node_type.#").HasValue("1"),
				check.That(data.ResourceName).Key("node_type.0.data_disk_size_gb").HasValue("140")),
		},
	})
}

func TestAccServiceFabricManagedCluster_authentication(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_service_fabric_managed_cluster", "test")
	r := ClusterResource{}
	nodeTypeData1 := r.nodeType("test1", true, 130)
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authentication(data, nodeTypeData1),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
	})
}

func (r ClusterResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := managedcluster.ParseManagedClusterID(state.ID)
	if err != nil {
		return nil, fmt.Errorf("while parsing resource ID: %+v", err)
	}

	resp, err := clients.ServiceFabricManaged.ManagedClusterClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("while checking for cluster's %q existence: %+v", id.String(), err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r ClusterResource) basic(data acceptance.TestData, nodeTypeData string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sfmc-%[1]d"
  location = "%[2]s"
}

resource "azurerm_service_fabric_managed_cluster" "test" {
  name                = "testacc-sfmc-%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard"
  username            = "testUser"
  password            = "NotV3ryS3cur3P@$$w0rd"
  dns_service_enabled = true

  client_connection_port = 12345
  http_gateway_port      = 23456

  lb_rule {
    backend_port       = 8000
    frontend_port      = 443
    probe_protocol     = "http"
    protocol           = "tcp"
    probe_request_path = "/"
  }

  %[4]s

  tags = {
    Test = "value"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, nodeTypeData)
}

func (r ClusterResource) withCustomSettings(data acceptance.TestData, nodeTypeData string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sfmc-%[1]d"
  location = "%[2]s"
}

resource "azurerm_service_fabric_managed_cluster" "test" {
  name                = "testacc-sfmc-%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard"
  username            = "testUser"
  password            = "NotV3ryS3cur3P@$$w0rd"
  dns_service_enabled = true

  client_connection_port = 12345
  http_gateway_port      = 23456

  lb_rule {
    backend_port       = 8000
    frontend_port      = 443
    probe_protocol     = "http"
    protocol           = "tcp"
    probe_request_path = "/"
  }

  custom_fabric_setting {
    section   = "ClusterManager"
    parameter = "EnableDefaultServicesUpgrade"
    value     = true
  }

  %[4]s

  tags = {
    Test = "value"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, nodeTypeData)
}

func (r ClusterResource) requiresImport(data acceptance.TestData, nt string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_service_fabric_managed_cluster" "test1" {
  name                = "testacc-sfmc-%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Standard"
  username            = "testUser"
  password            = "NotV3ryS3cur3P@$$w0rd"
  dns_service_enabled = true

  client_connection_port = 12345
  http_gateway_port      = 23456

  lb_rule {
    backend_port       = 8000
    frontend_port      = 443
    probe_protocol     = "http"
    protocol           = "tcp"
    probe_request_path = "/"
  }

  %[3]s

  tags = {
    Test = "value"
  }
}
`, r.basic(data, nt), data.RandomString, nt)
}

func (r ClusterResource) nodeType(name string, primary bool, diskSize int) string {
	return fmt.Sprintf(`
node_type {
  data_disk_size_gb      = %[1]d
  name                   = "%[2]s"
  primary                = %[3]t
  application_port_range = "7000-9000"
  ephemeral_port_range   = "10000-20000"

  vm_size            = "Standard_DS2_v2"
  vm_image_publisher = "MicrosoftWindowsServer"
  vm_image_sku       = "2016-Datacenter"
  vm_image_offer     = "WindowsServer"
  vm_image_version   = "latest"
  vm_instance_count  = 5
}
`, diskSize, name, primary)
}

func (r ClusterResource) authentication(data acceptance.TestData, nodeTypeData string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sfmc-%[1]d"
  location = "%[2]s"
}

resource "azurerm_service_fabric_managed_cluster" "test" {
  name                = "testacc-sfmc-%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  username            = "testUser"
  password            = "NotV3ryS3cur3P@$$w0rd"

  client_connection_port = 12345
  http_gateway_port      = 23456

  lb_rule {
    backend_port       = 8000
    frontend_port      = 443
    probe_protocol     = "http"
    protocol           = "tcp"
    probe_request_path = "/"
  }

  %[4]s

  authentication {
    certificate {
      thumbprint = "AAAA0982E0241795C04A61168D95B8DEE1B2CCCC"
      type       = "AdminClient"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, nodeTypeData)
}
