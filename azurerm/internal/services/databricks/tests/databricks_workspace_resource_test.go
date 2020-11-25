package tests

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databricks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databricks/parse"
)

func TestAzureRMDatabrickWorkspaceName(t *testing.T) {
	const errEmpty = "cannot be an empty string"
	const errMinLen = "must be at least 3 characters"
	const errMaxLen = "must be no more than 30 characters"
	const errAllowList = "can contain only alphanumeric characters, underscores, and hyphens"

	cases := []struct {
		Name           string
		Input          string
		ExpectedErrors []string
	}{
		// Happy paths:
		{
			Name:  "Entire character allow-list",
			Input: "aZ09_-",
		},
		{
			Name:  "Minimum character length",
			Input: "---",
		},
		{
			Name:  "Maximum character length",
			Input: "012345678901234567890123456789", // 30 chars
		},

		// Simple negative cases:
		{
			Name:           "Introduce a non-allowed character",
			Input:          "aZ09_-$", // dollar sign
			ExpectedErrors: []string{errAllowList},
		},
		{
			Name:           "Below minimum character length",
			Input:          "--",
			ExpectedErrors: []string{errMinLen},
		},
		{
			Name:           "Above maximum character length",
			Input:          "0123456789012345678901234567890", // 31 chars
			ExpectedErrors: []string{errMaxLen},
		},
		{
			Name:           "Specifically test for emptiness",
			Input:          "",
			ExpectedErrors: []string{errEmpty},
		},

		// Complex negative cases
		{
			Name:           "Too short and non-allowed char",
			Input:          "*^",
			ExpectedErrors: []string{errMinLen, errAllowList},
		},
		{
			Name:           "Too long and non-allowed char",
			Input:          "012345678901234567890123456789ß",
			ExpectedErrors: []string{errMaxLen, errAllowList},
		},
	}

	errsContain := func(errors []error, text string) bool {
		for _, err := range errors {
			if strings.Contains(err.Error(), text) {
				return true
			}
		}
		return false
	}

	t.Parallel()
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			_, errors := databricks.ValidateDatabricksWorkspaceName(tc.Input, "azurerm_databricks_workspace.test.name")

			if len(errors) != len(tc.ExpectedErrors) {
				t.Fatalf("Expected %d errors but got %d for %q: %v", len(tc.ExpectedErrors), len(errors), tc.Input, errors)
			}

			for _, expectedError := range tc.ExpectedErrors {
				if !errsContain(errors, expectedError) {
					t.Fatalf("Errors did not contain expected error: %s", expectedError)
				}
			}
		})
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
					resource.TestMatchResourceAttr(data.ResourceName, "workspace_url", regexp.MustCompile("azuredatabricks.net")),
					resource.TestCheckResourceAttrSet(data.ResourceName, "workspace_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDatabricksWorkspace_requiresImport(t *testing.T) {
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
			data.RequiresImportErrorStep(testAccAzureRMDatabricksWorkspace_requiresImport),
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
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Environment", "Production"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Pricing", "Standard"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDatabricksWorkspace_update(t *testing.T) {
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
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Environment", "Production"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Pricing", "Standard"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMDatabricksWorkspace_completeUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDatabricksWorkspaceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "managed_resource_group_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "managed_resource_group_name"),
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

		id, err := parse.WorkspaceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := conn.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return fmt.Errorf("Bad: Getting Workspace: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Databricks Workspace %s (resource group: %s) does not exist", id.Name, id.ResourceGroup)
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

		id, err := parse.WorkspaceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := conn.Get(ctx, id.ResourceGroup, id.Name)
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
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-db-%d"
  location = "%s"
}

resource "azurerm_databricks_workspace" "test" {
  name                = "acctestDBW-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "standard"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMDatabricksWorkspace_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMDatabricksWorkspace_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_databricks_workspace" "import" {
  name                = azurerm_databricks_workspace.test.name
  resource_group_name = azurerm_databricks_workspace.test.resource_group_name
  location            = azurerm_databricks_workspace.test.location
  sku                 = azurerm_databricks_workspace.test.sku
}
`, template)
}

func testAccAzureRMDatabricksWorkspace_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name = "acctestRG-db-%[1]d"

  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "public" {
  name                 = "acctest-sn-public-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
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
}

resource "azurerm_subnet" "private" {
  name                 = "acctest-sn-private-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
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
}

resource "azurerm_network_security_group" "nsg" {
  name                = "acctest-nsg-private-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet_network_security_group_association" "public" {
  subnet_id                 = azurerm_subnet.public.id
  network_security_group_id = azurerm_network_security_group.nsg.id
}

resource "azurerm_subnet_network_security_group_association" "private" {
  subnet_id                 = azurerm_subnet.private.id
  network_security_group_id = azurerm_network_security_group.nsg.id
}

resource "azurerm_databricks_workspace" "test" {
  name                        = "acctestDBW-%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku                         = "standard"
  managed_resource_group_name = "acctestRG-DBW-%[1]d-managed"

  custom_parameters {
    no_public_ip        = true
    public_subnet_name  = azurerm_subnet.public.name
    private_subnet_name = azurerm_subnet.private.name
    virtual_network_id  = azurerm_virtual_network.test.id
  }

  tags = {
    Environment = "Production"
    Pricing     = "Standard"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMDatabricksWorkspace_completeUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-db-%d"
  location = "%s"
}

resource "azurerm_databricks_workspace" "test" {
  name                        = "acctestDBW-%d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku                         = "standard"
  managed_resource_group_name = "acctestRG-DBW-%d-managed"

  tags = {
    Pricing = "Standard"
  }

  custom_parameters {
    no_public_ip = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
