package tests

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMAzureADServicePrincipal_byApplicationId(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_azuread_service_principal", "test")
	id := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryServicePrincipalDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMAzureADServicePrincipal_byApplicationId(id),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryServicePrincipalExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "application_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "object_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "display_name"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMAzureADServicePrincipal_byDisplayName(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_azuread_service_principal", "test")
	id := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryServicePrincipalDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMAzureADServicePrincipal_byDisplayName(id),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryServicePrincipalExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "application_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "object_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "display_name"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMAzureADServicePrincipal_byObjectId(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_azuread_service_principal", "test")
	id := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryServicePrincipalDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMAzureADServicePrincipal_byObjectId(id),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryServicePrincipalExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "application_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "object_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "display_name"),
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
