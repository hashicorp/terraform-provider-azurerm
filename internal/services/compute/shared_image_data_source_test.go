// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type SharedImageDataSource struct{}

func TestAccDataSourceSharedImage_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_shared_image", "test")
	r := SharedImageDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data, ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("hibernation_enabled").HasValue("false"),
			),
		},
	})
}

func TestAccDataSourceSharedImage_basicHyperVGenerationV2(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_shared_image", "test")
	r := SharedImageDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data, "V2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("hyper_v_generation").HasValue("V2"),
			),
		},
	})
}

func TestAccDataSourceSharedImage_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_shared_image", "test")
	r := SharedImageDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data, "V1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("hyper_v_generation").HasValue("V1"),
				check.That(data.ResourceName).Key("purchase_plan.0.name").HasValue("AccTestPlan"),
				check.That(data.ResourceName).Key("purchase_plan.0.publisher").HasValue("AccTestPlanPublisher"),
				check.That(data.ResourceName).Key("purchase_plan.0.product").HasValue("AccTestPlanProduct"),
			),
		},
	})
}

func TestAccDataSourceSharedImage_hibernationEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_shared_image", "test")
	r := SharedImageDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.withHibernationEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("hibernation_enabled").HasValue("true"),
			),
		},
	})
}

func TestAccDataSourceSharedImage_acceleratedNetworkSupportEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_shared_image", "test")
	r := SharedImageDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.withAcceleratedNetworkSupportEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("accelerated_network_support_enabled").HasValue("true"),
			),
		},
	})
}

func TestAccDataSourceSharedImage_trustedLaunchEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_shared_image", "test")
	r := SharedImageDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.withTrustedLaunchEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("trusted_launch_enabled").HasValue("true"),
			),
		},
	})
}

func TestAccDataSourceSharedImage_trustedLaunchSupported(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_shared_image", "test")
	r := SharedImageDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.withTrustedLaunchSupported(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("trusted_launch_supported").HasValue("true"),
			),
		},
	})
}

func TestAccDataSourceSharedImage_confidentialVMEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_shared_image", "test")
	r := SharedImageDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.withConfidentialVM(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("confidential_vm_enabled").HasValue("true"),
			),
		},
	})
}

func TestAccDataSourceSharedImage_confidentialVMSupported(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_shared_image", "test")
	r := SharedImageDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.withConfidentialVMSupported(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("confidential_vm_supported").HasValue("true"),
			),
		},
	})
}

func (SharedImageDataSource) basic(data acceptance.TestData, hyperVGen string) string {
	return fmt.Sprintf(`
%s

data "azurerm_shared_image" "test" {
  name                = azurerm_shared_image.test.name
  gallery_name        = azurerm_shared_image.test.gallery_name
  resource_group_name = azurerm_shared_image.test.resource_group_name
}
`, SharedImageResource{}.basicWithHyperVGen(data, hyperVGen))
}

func (SharedImageDataSource) complete(data acceptance.TestData, hyperVGen string) string {
	return fmt.Sprintf(`
%s

data "azurerm_shared_image" "test" {
  name                = azurerm_shared_image.test.name
  gallery_name        = azurerm_shared_image.test.gallery_name
  resource_group_name = azurerm_shared_image.test.resource_group_name
}
`, SharedImageResource{}.completeWithHyperVGen(data, hyperVGen))
}

func (SharedImageDataSource) withHibernationEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_shared_image" "test" {
  name                = azurerm_shared_image.test.name
  gallery_name        = azurerm_shared_image.test.gallery_name
  resource_group_name = azurerm_shared_image.test.resource_group_name
}
`, SharedImageResource{}.withHibernationEnabled(data))
}

func (SharedImageDataSource) withAcceleratedNetworkSupportEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_shared_image" "test" {
  name                = azurerm_shared_image.test.name
  gallery_name        = azurerm_shared_image.test.gallery_name
  resource_group_name = azurerm_shared_image.test.resource_group_name
}
`, SharedImageResource{}.withAcceleratedNetworkSupportEnabled(data))
}

func (SharedImageDataSource) withTrustedLaunchEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_shared_image" "test" {
  name                = azurerm_shared_image.test.name
  gallery_name        = azurerm_shared_image.test.gallery_name
  resource_group_name = azurerm_shared_image.test.resource_group_name
}
`, SharedImageResource{}.withTrustedLaunchEnabled(data))
}

func (SharedImageDataSource) withTrustedLaunchSupported(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_shared_image" "test" {
  name                = azurerm_shared_image.test.name
  gallery_name        = azurerm_shared_image.test.gallery_name
  resource_group_name = azurerm_shared_image.test.resource_group_name
}
`, SharedImageResource{}.withTrustedLaunchSupported(data))
}

func (SharedImageDataSource) withConfidentialVM(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_shared_image" "test" {
  name                = azurerm_shared_image.test.name
  gallery_name        = azurerm_shared_image.test.gallery_name
  resource_group_name = azurerm_shared_image.test.resource_group_name
}
`, SharedImageResource{}.withConfidentialVM(data))
}

func (SharedImageDataSource) withConfidentialVMSupported(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_shared_image" "test" {
  name                = azurerm_shared_image.test.name
  gallery_name        = azurerm_shared_image.test.gallery_name
  resource_group_name = azurerm_shared_image.test.resource_group_name
}
`, SharedImageResource{}.withConfidentialVmSupported(data))
}
