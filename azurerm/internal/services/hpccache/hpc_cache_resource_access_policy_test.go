package hpccache_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

func TestAccHPCCache_accessPolicy_default(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache", "test")
	r := HPCCacheResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.defaultAccessPolicyBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.defaultAccessPolicyComplete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.defaultAccessPolicyBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccHPCCache_accessPolicy_custom(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache", "test")
	r := HPCCacheResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.customAccessPolicySingleBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.customAccessPolicySingleComplete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.customAccessPolicyMultiple(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.customAccessPolicySingleBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r HPCCacheResource) defaultAccessPolicyBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache" "test" {
  name                = "acctest-HPCC-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cache_size_in_gb    = 3072
  subnet_id           = azurerm_subnet.test.id
  sku_name            = "Standard_2G"
  default_access_policy {
    access_rule {
      scope = "default"
      access = "rw"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r HPCCacheResource) defaultAccessPolicyComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache" "test" {
  name                = "acctest-HPCC-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cache_size_in_gb    = 3072
  subnet_id           = azurerm_subnet.test.id
  sku_name            = "Standard_2G"
  default_access_policy {
    access_rule {
      scope = "default"
      access = "ro"
    }

    access_rule {
      scope = "network"
      access = "rw"
	  filter = "10.0.0.0/24"
      suid_enabled = true
      submount_access_enabled = true
      root_squash_enabled = true
      anonymous_uid = 123
      anonymous_gid = 123
    }

    access_rule {
      scope = "host"
      access = "no"
	  filter = "10.0.0.1"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r HPCCacheResource) customAccessPolicySingleBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache" "test" {
  name                = "acctest-HPCC-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cache_size_in_gb    = 3072
  subnet_id           = azurerm_subnet.test.id
  sku_name            = "Standard_2G"
  custom_access_policy {
    name = "p1"
    access_rule {
      scope = "default"
      access = "rw"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r HPCCacheResource) customAccessPolicySingleComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache" "test" {
  name                = "acctest-HPCC-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cache_size_in_gb    = 3072
  subnet_id           = azurerm_subnet.test.id
  sku_name            = "Standard_2G"
  custom_access_policy {
    name = "p1"

    access_rule {
      scope = "default"
      access = "ro"
    }

    access_rule {
      scope = "network"
      access = "rw"
	  filter = "10.0.0.0/24"
      suid_enabled = true
      submount_access_enabled = true
      root_squash_enabled = true
      anonymous_uid = 123
      anonymous_gid = 123
    }

    access_rule {
      scope = "host"
      access = "no"
	  filter = "10.0.0.1"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r HPCCacheResource) customAccessPolicyMultiple(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache" "test" {
  name                = "acctest-HPCC-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cache_size_in_gb    = 3072
  subnet_id           = azurerm_subnet.test.id
  sku_name            = "Standard_2G"
  custom_access_policy {
    name = "p1"
    access_rule {
      scope = "default"
      access = "rw"
    }
  }

  custom_access_policy {
    name = "p2"
    access_rule {
      scope = "default"
      access = "ro"
    }
  }
}
`, r.template(data), data.RandomInteger)
}
