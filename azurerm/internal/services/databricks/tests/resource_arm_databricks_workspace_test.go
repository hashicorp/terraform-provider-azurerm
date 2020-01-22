package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databricks"
)

func TestAzureRMDatabrickWorkspaceName(t *testing.T) {
	cases := []struct {
		Value       string
		ShouldError bool
	}{
		{
			Value:       "hello",
			ShouldError: false,
		},
		{
			Value:       "hello123there",
			ShouldError: false,
		},
		{
			Value:       "hello-1-2-3-there",
			ShouldError: false,
		},
		{
			Value:       "hello-1-2-3-",
			ShouldError: true,
		},
		{
			Value:       "-hello-1-2-3",
			ShouldError: true,
		},
		{
			Value:       "hello!there",
			ShouldError: true,
		},
		{
			Value:       "hello--there",
			ShouldError: true,
		},
		{
			Value:       "!hellothere",
			ShouldError: true,
		},
		{
			Value:       "hellothere!",
			ShouldError: true,
		},
	}

	for _, tc := range cases {
		_, errors := databricks.ValidateDatabricksWorkspaceName(tc.Value, "test")

		hasErrors := len(errors) > 0
		if hasErrors && !tc.ShouldError {
			t.Fatalf("Expected no errors but got %d for %q", len(errors), tc.Value)
		}

		if !hasErrors && tc.ShouldError {
			t.Fatalf("Expected no errors but got %d for %q", len(errors), tc.Value)
		}
	}
}

func TestAccAzureRMDatabricksWorkspace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDatabricksWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDatabricksWorkspace_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDatabricksWorkspaceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "managed_resource_group_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDatabricksWorkspace_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDatabricksWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDatabricksWorkspace_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDatabricksWorkspaceExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMDatabricksWorkspace_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_databricks_workspace"),
			},
		},
	})
}

func TestAccAzureRMDatabricksWorkspace_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_workspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDatabricksWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDatabricksWorkspace_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDatabricksWorkspaceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "managed_resource_group_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "managed_resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "custom_parameters.0.virtual_network_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_parameters.0.public_subnet_name", "public"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_parameters.0.private_subnet_name", "private"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Environment", "Production"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Pricing", "Standard"),
				),
			},
			{
				Config: testAccAzureRMDatabricksWorkspace_completeUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDatabricksWorkspaceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "managed_resource_group_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "managed_resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "custom_parameters.0.virtual_network_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_parameters.0.public_subnet_name", "public"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_parameters.0.private_subnet_name", "private"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Pricing", "Standard"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMDatabricksWorkspaceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).DataBricks.WorkspacesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Bad: Not found: %s", resourceName)
		}

		workspaceName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: No resource group found in state for Databricks Workspace: %s", workspaceName)
		}

		resp, err := conn.Get(ctx, resourceGroup, workspaceName)
		if err != nil {
			return fmt.Errorf("Bad: Getting Workspace: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Databricks Workspace %s (resource group: %s) does not exist", workspaceName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDatabricksWorkspaceDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).DataBricks.WorkspacesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_databricks_workspace" {
			continue
		}

		workspaceName := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		resp, err := conn.Get(ctx, resourceGroup, workspaceName)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Bad: Databricks Workspace still exists:\n%#v", resp.ID)
		}
	}

	return nil
}

func testAccAzureRMDatabricksWorkspace_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_databricks_workspace" "test" {
  name                = "acctestdbw-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  sku                 = "standard"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMDatabricksWorkspace_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMDatabricksWorkspace_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_databricks_workspace" "import" {
  name                = "$[azurerm_databricks_workspace.test.name}"
  resource_group_name = "${azurerm_databricks_workspace.test.resource_group_name}"
  location            = "${azurerm_databricks_workspace.test.location}"
  sku                 = "${azurerm_databricks_workspace.test.sku}"
}
`, template)
}

func testAccAzureRMDatabricksWorkspace_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "test"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "public" {
  name                 = "public"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.1.0/24"

  delegation {
    name = "acctest"

    service_delegation {
      name = "Microsoft.Databricks/workspaces"

      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
        "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
      ]
    }
  }

  lifecycle {
    ignore_changes = ["network_security_group_id"]
  }
}

resource "azurerm_subnet" "private" {
  name                 = "private"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"

  delegation {
    name = "acctest"

    service_delegation {
      name = "Microsoft.Databricks/workspaces"

      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
        "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
      ]
    }
  }

  lifecycle {
    ignore_changes = ["network_security_group_id"]
  }
}

resource "azurerm_network_security_group" "nsg" {
  name                = "private-nsg"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet_network_security_group_association" "public" {
  subnet_id                 = "${azurerm_subnet.public.id}"
  network_security_group_id = "${azurerm_network_security_group.nsg.id}"
}

resource "azurerm_subnet_network_security_group_association" "private" {
  subnet_id                 = "${azurerm_subnet.private.id}"
  network_security_group_id = "${azurerm_network_security_group.nsg.id}"
}

resource "azurerm_databricks_workspace" "test" {
  name                        = "acctestdbw-%d"
  resource_group_name         = "${azurerm_resource_group.test.name}"
  location                    = "${azurerm_resource_group.test.location}"
  sku                         = "standard"
  managed_resource_group_name = "acctestRG-%d-managed"

  custom_parameters {
    public_subnet_name  = "${azurerm_subnet.public.name}"
    private_subnet_name = "${azurerm_subnet.private.name}"
    virtual_network_id  = "${azurerm_virtual_network.test.id}"
  }

  tags = {
    Environment = "Production"
    Pricing     = "Standard"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDatabricksWorkspace_completeUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_databricks_workspace" "test" {
  name                        = "acctestdbw-%d"
  resource_group_name         = "${azurerm_resource_group.test.name}"
  location                    = "${azurerm_resource_group.test.location}"
  sku                         = "standard"
  managed_resource_group_name = "acctestRG-%d-managed"

  tags = {
    Pricing = "Standard"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
