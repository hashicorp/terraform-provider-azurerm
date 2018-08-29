package azurerm

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMAzureADGroup_byObjectId(t *testing.T) {
	dataSourceName := "data.azurerm_azuread_group.test"
	id := uuid.New().String()
	config := testAccDataSourceAzureRMAzureADGroup_objectId(id)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMActiveDirectoryGroup(id),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryGroupExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "name", fmt.Sprintf("acctest%s", id)),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMAzureADGroup_byName(t *testing.T) {
	dataSourceName := "data.azurerm_azuread_group.test"
	id := uuid.New().String()
	config := testAccDataSourceAzureRMAzureADGroup_name(id)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMActiveDirectoryGroup(id),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryGroupExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "name", fmt.Sprintf("acctest%s", id)),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMAzureADGroup_objectId(id string) string {
	template := testAccAzureRMActiveDirectoryGroup(id)
	return fmt.Sprintf(`
%s

data "azurerm_azuread_group" "test" {
  object_id = "${azurerm_azuread_group.test.id}"
}
`, template)
}

func testAccDataSourceAzureRMAzureADGroup_name(id string) string {
	template := testAccAzureRMActiveDirectoryGroup(id)
	return fmt.Sprintf(`
%s

data "azurerm_azuread_group" "test" {
  name = "${azurerm_azuread_group.test.name}"
}
`, template)
}
