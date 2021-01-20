package tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type KubernetesServiceVersionDataSource struct {
}

const k8sVersionRX = `[0-9]+\.[0-9]+\.[0-9]*`

func TestAccDataSourceAzureRMKubernetesServiceVersions_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_service_versions", "test")
	r := KubernetesServiceVersionDataSource{}
	kvrx := regexp.MustCompile(k8sVersionRX)

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("versions.#").Exists(),
				resource.TestMatchResourceAttr(data.ResourceName, "versions.0", kvrx),
				check.That(data.ResourceName).Key("latest_version").Exists(),
				resource.TestMatchResourceAttr(data.ResourceName, "latest_version", kvrx),
			),
		},
	})
}

func TestAccDataSourceAzureRMKubernetesServiceVersions_filtered(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_service_versions", "test")
	r := KubernetesServiceVersionDataSource{}
	kvrx := regexp.MustCompile(k8sVersionRX)

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.filtered(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("versions.#").Exists(),
				resource.TestMatchResourceAttr(data.ResourceName, "versions.0", kvrx),
				check.That(data.ResourceName).Key("latest_version").Exists(),
				resource.TestMatchResourceAttr(data.ResourceName, "latest_version", kvrx),
			),
		},
	})
}

func TestAccDataSourceAzureRMKubernetesServiceVersions_nopreview(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_service_versions", "test")
	r := KubernetesServiceVersionDataSource{}
	kvrx := regexp.MustCompile(k8sVersionRX)

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.nopreview(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("versions.#").Exists(),
				resource.TestMatchResourceAttr(data.ResourceName, "versions.0", kvrx),
				check.That(data.ResourceName).Key("latest_version").Exists(),
				resource.TestMatchResourceAttr(data.ResourceName, "latest_version", kvrx),
			),
		},
	})
}

func (KubernetesServiceVersionDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_kubernetes_service_versions" "test" {
  location = "%s"
}
`, data.Locations.Primary)
}

func (KubernetesServiceVersionDataSource) filtered(data acceptance.TestData) string {
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

func (KubernetesServiceVersionDataSource) nopreview(data acceptance.TestData) string {
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
