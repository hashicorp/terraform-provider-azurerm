package azurerm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

const k8sVersionRX = `[0-9]+\.[0-9]+\.[0-9]*`

func TestAccDataSourceAzureRMKubernetesServiceVersions_basic(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_service_versions.test"
	location := acceptance.Location()
	kvrx := regexp.MustCompile(k8sVersionRX)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesServiceVersions_basic(location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "versions.#"),
					resource.TestMatchResourceAttr(dataSourceName, "versions.0", kvrx),
					resource.TestCheckResourceAttrSet(dataSourceName, "latest_version"),
					resource.TestMatchResourceAttr(dataSourceName, "latest_version", kvrx),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKubernetesServiceVersions_filtered(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_service_versions.test"
	location := acceptance.Location()
	kvrx := regexp.MustCompile(k8sVersionRX)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesServiceVersions_filtered(location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "versions.#"),
					resource.TestMatchResourceAttr(dataSourceName, "versions.0", kvrx),
					resource.TestCheckResourceAttrSet(dataSourceName, "latest_version"),
					resource.TestMatchResourceAttr(dataSourceName, "latest_version", kvrx),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMKubernetesServiceVersions_basic(location string) string {
	return fmt.Sprintf(`
data "azurerm_kubernetes_service_versions" "test" {
  location = "%s"
}
`, location)
}

func testAccDataSourceAzureRMKubernetesServiceVersions_filtered(location string) string {
	return fmt.Sprintf(`
data "azurerm_kubernetes_service_versions" "test" {
  location       = "%s"
  version_prefix = "1."
}
`, location)
}
