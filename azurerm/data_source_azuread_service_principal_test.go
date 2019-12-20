package azurerm

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMAzureADServicePrincipal_byApplicationId(t *testing.T) {
	dataSourceName := "data.azurerm_azuread_service_principal.test"
	id := uuid.New().String()
	config := testAccDataSourceAzureRMAzureADServicePrincipal_byApplicationId(id)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryServicePrincipalDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryServicePrincipalExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "application_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "object_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "display_name"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMAzureADServicePrincipal_byDisplayName(t *testing.T) {
	dataSourceName := "data.azurerm_azuread_service_principal.test"
	id := uuid.New().String()
	config := testAccDataSourceAzureRMAzureADServicePrincipal_byDisplayName(id)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryServicePrincipalDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryServicePrincipalExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "application_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "object_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "display_name"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMAzureADServicePrincipal_byObjectId(t *testing.T) {
	dataSourceName := "data.azurerm_azuread_service_principal.test"
	id := uuid.New().String()
	config := testAccDataSourceAzureRMAzureADServicePrincipal_byObjectId(id)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryServicePrincipalDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryServicePrincipalExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "application_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "object_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "display_name"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMAzureADServicePrincipal_byApplicationId(id string) string {
	template := testAccAzureRMActiveDirectoryServicePrincipal_basic(id)
	return fmt.Sprintf(`
%s

data "azurerm_azuread_service_principal" "test" {
  application_id = "${azurerm_azuread_service_principal.test.application_id}"
}
`, template)
}

func testAccDataSourceAzureRMAzureADServicePrincipal_byDisplayName(id string) string {
	template := testAccAzureRMActiveDirectoryServicePrincipal_basic(id)
	return fmt.Sprintf(`
%s

data "azurerm_azuread_service_principal" "test" {
  display_name = "${azurerm_azuread_service_principal.test.display_name}"
}
`, template)
}

func testAccDataSourceAzureRMAzureADServicePrincipal_byObjectId(id string) string {
	template := testAccAzureRMActiveDirectoryServicePrincipal_basic(id)
	return fmt.Sprintf(`
%s

data "azurerm_azuread_service_principal" "test" {
  object_id = "${azurerm_azuread_service_principal.test.id}"
}
`, template)
}
