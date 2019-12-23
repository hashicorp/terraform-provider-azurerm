package azurerm

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMAzureADApplication_byObjectId(t *testing.T) {
	dataSourceName := "data.azurerm_azuread_application.test"
	id := uuid.New().String()
	config := testAccDataSourceAzureRMAzureADApplication_objectId(id)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMActiveDirectoryApplication_basic(id),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryApplicationExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "name", fmt.Sprintf("acctest%s", id)),
					resource.TestCheckResourceAttr(dataSourceName, "homepage", fmt.Sprintf("https://acctest%s", id)),
					resource.TestCheckResourceAttr(dataSourceName, "identifier_uris.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "reply_urls.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "oauth2_allow_implicit_flow", "false"),
					resource.TestCheckResourceAttrSet(dataSourceName, "application_id"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMAzureADApplication_byObjectIdComplete(t *testing.T) {
	dataSourceName := "data.azurerm_azuread_application.test"
	id := uuid.New().String()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMActiveDirectoryApplication_complete(id),
			},
			{
				Config: testAccDataSourceAzureRMAzureADApplication_objectIdComplete(id),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryApplicationExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "name", fmt.Sprintf("acctest%s", id)),
					resource.TestCheckResourceAttr(dataSourceName, "homepage", fmt.Sprintf("https://homepage-%s", id)),
					resource.TestCheckResourceAttr(dataSourceName, "identifier_uris.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "reply_urls.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "oauth2_allow_implicit_flow", "true"),
					resource.TestCheckResourceAttrSet(dataSourceName, "application_id"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMAzureADApplication_byName(t *testing.T) {
	dataSourceName := "data.azurerm_azuread_application.test"
	id := uuid.New().String()
	config := testAccDataSourceAzureRMAzureADApplication_name(id)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMActiveDirectoryApplication_basic(id),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryApplicationExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "name", fmt.Sprintf("acctest%s", id)),
					resource.TestCheckResourceAttr(dataSourceName, "homepage", fmt.Sprintf("https://acctest%s", id)),
					resource.TestCheckResourceAttr(dataSourceName, "identifier_uris.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "reply_urls.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "oauth2_allow_implicit_flow", "false"),
					resource.TestCheckResourceAttrSet(dataSourceName, "application_id"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMAzureADApplication_objectId(id string) string {
	template := testAccAzureRMActiveDirectoryApplication_basic(id)
	return fmt.Sprintf(`
%s

data "azurerm_azuread_application" "test" {
  object_id = "${azurerm_azuread_application.test.id}"
}
`, template)
}

func testAccDataSourceAzureRMAzureADApplication_objectIdComplete(id string) string {
	template := testAccAzureRMActiveDirectoryApplication_complete(id)
	return fmt.Sprintf(`
%s

data "azurerm_azuread_application" "test" {
  object_id = "${azurerm_azuread_application.test.id}"
}
`, template)
}

func testAccDataSourceAzureRMAzureADApplication_name(id string) string {
	template := testAccAzureRMActiveDirectoryApplication_basic(id)
	return fmt.Sprintf(`
%s

data "azurerm_azuread_application" "test" {
  name = "${azurerm_azuread_application.test.name}"
}
`, template)
}
