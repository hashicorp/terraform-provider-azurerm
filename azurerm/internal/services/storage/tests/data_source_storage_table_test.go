package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureStorageTable_resourceId(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_table", "resource_id")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureStorageTable_resourceId(data),
			},
			{
				Config: testAccDataSourceAzureStorageTable_resourceIdWithDataSource(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckTableEntity2DataSource(data.ResourceName, "azurerm_storage_table_entity.resource_id"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureStorageTable_key(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_table", "key")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureStorageTable_key(data),
			},
			{
				Config: testAccDataSourceAzureStorageTable_keyWithDataSource(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckTableEntity2DataSource(data.ResourceName, "azurerm_storage_table_entity.key"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureStorageTable_query(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_storage_table", "query")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureStorageTable_query(data),
			},
			{
				Config: testAccDataSourceAzureStorageTable_queryWithDataSource(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckTableEntity2DataSource(data.ResourceName, "azurerm_storage_table_entity.query1"),
					testCheckTableEntity2DataSource(data.ResourceName, "azurerm_storage_table_entity.query2"),
				),
			},
		},
	})
}

func getResourceLabel(data acceptance.TestData) string {
	resourceNameParts := strings.Split(data.ResourceName, ".")
	resourceLabel := resourceNameParts[len(resourceNameParts)-1]

	return resourceLabel
}

func getResourceDefinitionsCommon(data acceptance.TestData) string {
	resourceLabel := getResourceLabel(data)

	return fmt.Sprintf(`
resource "azurerm_resource_group" "%s" {
  name     = "acctest-dsast-%s"
  location = "%s"
}

resource "azurerm_storage_account" "%s" {
  name                     = "acctestdsast%s"
  resource_group_name      = "${azurerm_resource_group.%s.name}"
  location                 = "${azurerm_resource_group.%s.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "test"
  }
}

resource "azurerm_storage_table" "%s" {
  name                 = "acctestdsast%s"
  storage_account_name = "${azurerm_storage_account.%s.name}"
}
`, resourceLabel, data.RandomString, data.Locations.Primary, resourceLabel, data.RandomString, resourceLabel, resourceLabel, resourceLabel, data.RandomString, resourceLabel)
}

func testAccDataSourceAzureStorageTable_resourceId(data acceptance.TestData) string {
	resourceDefinitionsCommon := getResourceDefinitionsCommon(data)
	resourceLabel := getResourceLabel(data)

	return fmt.Sprintf(`
%s
resource "azurerm_storage_table_entity" "%s" {
  storage_account_name = "${azurerm_storage_account.%s.name}"
  table_name           = "${azurerm_storage_table.%s.name}"
  partition_key        = "mypartition"
  row_key              = "1"

  entity = {
    breed     = "Scottish Fold"
    life_span = "11 - 14"
    img_url   = "https://cdn2.thecatapi.com/images/4UTnt4f74.jpg"
  }
}
`, resourceDefinitionsCommon, resourceLabel, resourceLabel, resourceLabel)
}

func testAccDataSourceAzureStorageTable_resourceIdWithDataSource(data acceptance.TestData) string {
	config := testAccDataSourceAzureStorageTable_resourceId(data)
	resourceLabel := getResourceLabel(data)

	return fmt.Sprintf(`
%s
data "azurerm_storage_table" "%s" {
  resource_id = azurerm_storage_table_entity.%s.id
}
`, config, resourceLabel, resourceLabel)
}

func testAccDataSourceAzureStorageTable_key(data acceptance.TestData) string {
	resourceDefinitionsCommon := getResourceDefinitionsCommon(data)
	resourceLabel := getResourceLabel(data)

	return fmt.Sprintf(`
%s
resource "azurerm_storage_table_entity" "%s" {
  storage_account_name = "${azurerm_storage_account.%s.name}"
  table_name           = "${azurerm_storage_table.%s.name}"
  partition_key        = "mypartition"
  row_key              = "2"

  entity = {
    breed     = "Munchkin"
    life_span = "10 - 15"
    img_url   = "https://cdn2.thecatapi.com/images/hxlto6Z4I.jpg"
  }
}
`, resourceDefinitionsCommon, resourceLabel, resourceLabel, resourceLabel)
}

func testAccDataSourceAzureStorageTable_keyWithDataSource(data acceptance.TestData) string {
	config := testAccDataSourceAzureStorageTable_key(data)
	resourceLabel := getResourceLabel(data)

	return fmt.Sprintf(`
%s
data "azurerm_storage_table" "%s" {
  key {
    storage_account_name = "${azurerm_storage_account.%s.name}"
    table_name           = "${azurerm_storage_table.%s.name}"
    partition_key        = "mypartition"
    row_key              = "2"
  }
}
`, config, resourceLabel, resourceLabel, resourceLabel)
}

func testAccDataSourceAzureStorageTable_query(data acceptance.TestData) string {
	resourceDefinitionsCommon := getResourceDefinitionsCommon(data)
	resourceLabel := getResourceLabel(data)

	return fmt.Sprintf(`
%s
resource "azurerm_storage_table_entity" "%s1" {
  storage_account_name = "${azurerm_storage_account.%s.name}"
  table_name           = "${azurerm_storage_table.%s.name}"
  partition_key        = "mypartition"
  row_key              = "1"

  entity = {
    breed     = "Scottish Fold"
    life_span = "11 - 14"
    img_url   = "https://cdn2.thecatapi.com/images/4UTnt4f74.jpg"
  }
}

resource "azurerm_storage_table_entity" "%s2" {
  storage_account_name = "${azurerm_storage_account.%s.name}"
  table_name           = "${azurerm_storage_table.%s.name}"
  partition_key        = "mypartition"
  row_key              = "2"

  entity = {
    breed     = "Munchkin"
    life_span = "10 - 15"
    img_url   = "https://cdn2.thecatapi.com/images/hxlto6Z4I.jpg"
  }
}
`, resourceDefinitionsCommon, resourceLabel, resourceLabel, resourceLabel, resourceLabel, resourceLabel, resourceLabel)
}

func testAccDataSourceAzureStorageTable_queryWithDataSource(data acceptance.TestData) string {
	config := testAccDataSourceAzureStorageTable_query(data)
	resourceLabel := getResourceLabel(data)

	return fmt.Sprintf(`
%s
data "azurerm_storage_table" "%s" {
  query {
    storage_account_name = "${azurerm_storage_account.%s.name}"
    table_name           = "${azurerm_storage_table.%s.name}"
    filter               = "PartitionKey eq 'mypartition'"
  }
}
`, config, resourceLabel, resourceLabel, resourceLabel)
}

func testCheckTableEntity2DataSource(dataSourceResourceName string, tableEntityResourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		dataSource := s.RootModule().Resources[dataSourceResourceName]
		if dataSource == nil {
			return fmt.Errorf("Cannot find data-source resource in TF state")
		}

		tableEntity := s.RootModule().Resources[tableEntityResourceName]
		if tableEntity == nil {
			return fmt.Errorf("Cannot find Table entity resource in TF state")
		}

		// go through tableEntity "entity.*"
		for kTe, vTe := range tableEntity.Primary.Attributes {
			if strings.Contains(kTe, "entity.") {
				// extract <key> from "entity.<key>"
				kParts := strings.Split(kTe, ".")
				kPartsLast := kParts[len(kParts)-1]
				wasKeyFound := false
				wasKeyEqual := false
				// go through data-source
				for kDs, vDs := range dataSource.Primary.Attributes {
					// find  <key> in data-source attrs
					if strings.Contains(kDs, "."+kPartsLast) {
						wasKeyFound = true

						if vTe == vDs {
							wasKeyEqual = true

							break
						}
					}
				}
				fmt.Println("")

				if wasKeyFound == false {
					return fmt.Errorf("Checking datasource failed - key %q in Table entity not found in data-source", kTe)
				}
				if wasKeyEqual == false {
					return fmt.Errorf("Checking datasource failed - unable to find matching field for %q:%q in data-source data", kTe, vTe)
				}
			}
		}

		return nil
	}
}
