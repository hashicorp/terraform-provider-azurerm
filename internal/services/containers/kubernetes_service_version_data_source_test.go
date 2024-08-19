// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type KubernetesServiceVersionDataSource struct{}

const k8sVersionRX = `[0-9]+\.[0-9]+\.[0-9]*`

func TestAccDataSourceAzureRMKubernetesServiceVersions_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_service_versions", "test")
	r := KubernetesServiceVersionDataSource{}
	kvrx := regexp.MustCompile(k8sVersionRX)

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("versions.#").Exists(),
				acceptance.TestMatchResourceAttr(data.ResourceName, "versions.0", kvrx),
				check.That(data.ResourceName).Key("latest_version").Exists(),
				acceptance.TestMatchResourceAttr(data.ResourceName, "latest_version", kvrx),
			),
		},
	})
}

func TestAccDataSourceAzureRMKubernetesServiceVersions_filtered(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_service_versions", "test")
	r := KubernetesServiceVersionDataSource{}
	kvrx := regexp.MustCompile(k8sVersionRX)

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.filtered(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("versions.#").Exists(),
				acceptance.TestMatchResourceAttr(data.ResourceName, "versions.0", kvrx),
				check.That(data.ResourceName).Key("latest_version").Exists(),
				acceptance.TestMatchResourceAttr(data.ResourceName, "latest_version", kvrx),
				check.That(data.ResourceName).Key("default_version").Exists(),
			),
		},
	})
}

func TestAccDataSourceAzureRMKubernetesServiceVersions_nopreview(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_service_versions", "test")
	r := KubernetesServiceVersionDataSource{}
	kvrx := regexp.MustCompile(k8sVersionRX)

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.nopreview(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("versions.#").Exists(),
				acceptance.TestMatchResourceAttr(data.ResourceName, "versions.0", kvrx),
				check.That(data.ResourceName).Key("latest_version").Exists(),
				acceptance.TestMatchResourceAttr(data.ResourceName, "latest_version", kvrx),
				check.That(data.ResourceName).Key("default_version").Exists(),
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
