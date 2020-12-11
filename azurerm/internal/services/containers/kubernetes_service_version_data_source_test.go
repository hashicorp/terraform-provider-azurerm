package containers_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

const k8sVersionRX = `[0-9]+\.[0-9]+\.[0-9]*`

func TestAccDataSourceAzureRMKubernetesServiceVersions_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_service_versions", "test")
	kvrx := regexp.MustCompile(k8sVersionRX)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesServiceVersions_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "versions.#"),
					resource.TestMatchResourceAttr(data.ResourceName, "versions.0", kvrx),
					resource.TestCheckResourceAttrSet(data.ResourceName, "latest_version"),
					resource.TestMatchResourceAttr(data.ResourceName, "latest_version", kvrx),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKubernetesServiceVersions_filtered(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_service_versions", "test")
	kvrx := regexp.MustCompile(k8sVersionRX)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesServiceVersions_filtered(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "versions.#"),
					resource.TestMatchResourceAttr(data.ResourceName, "versions.0", kvrx),
					resource.TestCheckResourceAttrSet(data.ResourceName, "latest_version"),
					resource.TestMatchResourceAttr(data.ResourceName, "latest_version", kvrx),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKubernetesServiceVersions_nopreview(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_service_versions", "test")
	kvrx := regexp.MustCompile(k8sVersionRX)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesServiceVersions_nopreview(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "versions.#"),
					resource.TestMatchResourceAttr(data.ResourceName, "versions.0", kvrx),
					resource.TestCheckResourceAttrSet(data.ResourceName, "latest_version"),
					resource.TestMatchResourceAttr(data.ResourceName, "latest_version", kvrx),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMKubernetesServiceVersions_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_kubernetes_service_versions" "test" {
  location = "%s"
}
`, data.Locations.Primary)
}

func testAccDataSourceAzureRMKubernetesServiceVersions_filtered(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_kubernetes_service_versions" "test" {
  location       = "%s"
  version_prefix = "1."
}
`, data.Locations.Primary)
}

func testAccDataSourceAzureRMKubernetesServiceVersions_nopreview(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_kubernetes_service_versions" "test" {
  location        = "%s"
  include_preview = false
}
`, data.Locations.Primary)
}
