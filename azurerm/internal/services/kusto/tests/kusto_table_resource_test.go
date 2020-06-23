package tests

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-kusto-go/kusto"
	"github.com/Azure/azure-kusto-go/kusto/data/table"
	"github.com/Azure/azure-kusto-go/kusto/unsafe"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/kusto/parse"
	dataplaneTypes "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/kusto/types"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMKustoTable_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_table", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKustoTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKustoTable_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKustoTableExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "database_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "column.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "column.0.name", "col1"),
					resource.TestCheckResourceAttr(data.ResourceName, "column.0.type", "string"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKustoTable_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_table", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKustoTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKustoTable_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKustoTableExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "database_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "folder", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "doc", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "column.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "column.0.name", "col1"),
					resource.TestCheckResourceAttr(data.ResourceName, "column.0.type", "string"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMKustoTable_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKustoTableExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "database_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "folder", "folder"),
					resource.TestCheckResourceAttr(data.ResourceName, "doc", "documentation"),
					resource.TestCheckResourceAttr(data.ResourceName, "column.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "column.0.name", "col2"),
					resource.TestCheckResourceAttr(data.ResourceName, "column.0.type", "real"),
					resource.TestCheckResourceAttr(data.ResourceName, "column.1.name", "col1"),
					resource.TestCheckResourceAttr(data.ResourceName, "column.1.type", "string"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMKustoTable_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKustoTableExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "database_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "folder", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "doc", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "column.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "column.0.name", "col1"),
					resource.TestCheckResourceAttr(data.ResourceName, "column.0.type", "string"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKustoTable_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_table", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKustoTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKustoTable_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKustoTableExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMKustoTable_requiresImport),
		},
	})
}

func TestAccAzureRMKustoTable_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_table", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKustoTable_complete(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "database_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "folder", "folder"),
					resource.TestCheckResourceAttr(data.ResourceName, "doc", "documentation"),
					resource.TestCheckResourceAttr(data.ResourceName, "column.#", "10"),
					resource.TestCheckResourceAttr(data.ResourceName, "column.0.name", "col1"),
					resource.TestCheckResourceAttr(data.ResourceName, "column.0.type", "bool"),
					resource.TestCheckResourceAttr(data.ResourceName, "column.1.name", "col2"),
					resource.TestCheckResourceAttr(data.ResourceName, "column.1.type", "datetime"),
					resource.TestCheckResourceAttr(data.ResourceName, "column.2.name", "col3"),
					resource.TestCheckResourceAttr(data.ResourceName, "column.2.type", "decimal"),
					resource.TestCheckResourceAttr(data.ResourceName, "column.3.name", "col4"),
					resource.TestCheckResourceAttr(data.ResourceName, "column.3.type", "dynamic"),
					resource.TestCheckResourceAttr(data.ResourceName, "column.4.name", "col5"),
					resource.TestCheckResourceAttr(data.ResourceName, "column.4.type", "guid"),
					resource.TestCheckResourceAttr(data.ResourceName, "column.5.name", "col6"),
					resource.TestCheckResourceAttr(data.ResourceName, "column.5.type", "int"),
					resource.TestCheckResourceAttr(data.ResourceName, "column.6.name", "col7"),
					resource.TestCheckResourceAttr(data.ResourceName, "column.6.type", "long"),
					resource.TestCheckResourceAttr(data.ResourceName, "column.7.name", "col8"),
					resource.TestCheckResourceAttr(data.ResourceName, "column.7.type", "real"),
					resource.TestCheckResourceAttr(data.ResourceName, "column.8.name", "col9"),
					resource.TestCheckResourceAttr(data.ResourceName, "column.8.type", "string"),
					resource.TestCheckResourceAttr(data.ResourceName, "column.9.name", "col10"),
					resource.TestCheckResourceAttr(data.ResourceName, "column.9.type", "timespan"),
				),
			},
		},
	})
}

func testAccAzureRMKustoTable_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_kusto_cluster" "cluster" {
  name                = "acctestkc%[3]s"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}

resource "azurerm_kusto_database" "test" {
  name                = "acctestkd-%[1]d"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  cluster_name        = azurerm_kusto_cluster.cluster.name
}

resource "azurerm_kusto_table" "test" {
  name        = "acctestkt%[1]d"
  database_id = azurerm_kusto_database.test.id
  column {
    name = "col1"
    type = "string"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMKustoTable_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMKustoTable_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_kusto_table" "import" {
  name        = azurerm_kusto_table.test.name
  database_id = azurerm_kusto_database.test.id
  column {
    name = azurerm_kusto_table.test.column.0.name
    type = azurerm_kusto_table.test.column.0.type
  }
}
`, template)
}

func testAccAzureRMKustoTable_update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_kusto_cluster" "cluster" {
  name                = "acctestkc%[3]s"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}

resource "azurerm_kusto_database" "test" {
  name                = "acctestkd-%[1]d"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  cluster_name        = azurerm_kusto_cluster.cluster.name
}

resource "azurerm_kusto_table" "test" {
  name        = "acctestkt%[1]d"
  database_id = azurerm_kusto_database.test.id
  folder      = "folder"
  doc         = "documentation"
  column {
    name = "col2"
    type = "real"
  }
  column {
    name = "col1"
    type = "string"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMKustoTable_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_kusto_cluster" "cluster" {
  name                = "acctestkc%[3]s"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}

resource "azurerm_kusto_database" "test" {
  name                = "acctestkd-%[1]d"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  cluster_name        = azurerm_kusto_cluster.cluster.name
}

resource "azurerm_kusto_table" "test" {
  name        = "acctestkt%[1]d"
  database_id = azurerm_kusto_database.test.id
  doc         = "documentation"
  folder      = "folder"

  column {
    name = "col1"
    type = "bool"
  }
  column {
    name = "col2"
    type = "datetime"
  }
  column {
    name = "col3"
    type = "decimal"
  }
  column {
    name = "col4"
    type = "dynamic"
  }
  column {
    name = "col5"
    type = "guid"
  }
  column {
    name = "col6"
    type = "int"
  }
  column {
    name = "col7"
    type = "long"
  }
  column {
    name = "col8"
    type = "real"
  }
  column {
    name = "col9"
    type = "string"
  }
  column {
    name = "col10"
    type = "timespan"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testCheckAzureRMKustoTableDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_kusto_table" {
			continue
		}

		databaseID := rs.Primary.Attributes["database_id"]
		name := rs.Primary.Attributes["name"]

		id, err := parse.KustoDatabaseID(databaseID)
		if err != nil {
			return err
		}

		exists, err := checkTableExists(rs, id, name)

		if err != nil {
			return err
		}

		if exists == nil || *exists == false {
			return nil
		}

		return fmt.Errorf("Bad: Kusto Table %q (resource group: %q, cluster: %q, database: %q) does still exist", name, id.ResourceGroup, id.Cluster, id.Name)
	}

	return nil
}

func testCheckAzureRMKustoTableExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		databaseID := rs.Primary.Attributes["database_id"]

		id, err := parse.KustoDatabaseID(databaseID)
		if err != nil {
			return err
		}

		exists, err := checkTableExists(rs, id, name)

		if err != nil {
			return fmt.Errorf("Bad: Get on Kusto Data Plane Client: %+v", err)
		}

		if exists == nil || *exists == false {
			return fmt.Errorf("Bad: Kusto Table %q (resource group: %q, cluster: %q, database: %q) does not exist", name, id.ResourceGroup, id.Cluster, id.Name)
		}

		return nil
	}
}

func checkTableExists(rs *terraform.ResourceState, id *parse.KustoDatabaseId, name string) (*bool, error) {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Kusto.ClustersClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	cluster, err := client.Get(ctx, id.ResourceGroup, id.Cluster)
	if err != nil {
		if utils.ResponseWasNotFound(cluster.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Kusto Cluster %q (Resource Group %q): %+v", id.Cluster, id.ResourceGroup, err)
	}
	if cluster.ClusterProperties == nil || cluster.ClusterProperties.URI == nil {
		return nil, fmt.Errorf("Kusto Cluster %q (Resource Group %q) URI property is nil or empty", id.Cluster, id.ResourceGroup)
	}

	dataplaneClient, err := acceptance.AzureProvider.Meta().(*clients.Client).Kusto.NewDataPlaneClient(*cluster.URI)
	if err != nil {
		return nil, fmt.Errorf("init Kusto Data Plane Client: %+v", err)
	}

	stmtRaw := fmt.Sprintf(".show tables | where TableName == \"%s\"", name)
	stmt := kusto.NewStmt("", kusto.UnsafeStmt(unsafe.Stmt{Add: true})).UnsafeAdd(stmtRaw)
	iter, err := dataplaneClient.Mgmt(ctx, id.Name, stmt)
	if err != nil {
		return nil, fmt.Errorf("querying Kusto Table %q (Cluster %q, Database %q): %+v", name, id.Cluster, id.Name, err)
	}
	defer iter.Stop()

	found := false
	err = iter.Do(
		func(row *table.Row) error {
			rec := dataplaneTypes.KustoTableRecord{}
			if err := row.ToStruct(&rec); err != nil {
				return err
			}
			found = true
			return nil
		},
	)

	return &found, nil
}
