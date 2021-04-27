package hpccache_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/hpccache/parse"
	networkParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type HPCCacheResource struct {
}

type HPCCacheDirectoryADInfo struct {
	SubnetID networkParse.SubnetId

	PrimaryDNS        string
	DomainName        string
	CacheNetBiosName  string
	DomainNetBiosName string
	Username          string
	Password          string
}

const (
	adSubnetIDEnv          = "ARM_TEST_HPC_AD_SUBNET_ID"
	adPrimaryDNSEnv        = "ARM_TEST_HPC_AD_PRIMARY_DNS"
	adDomainNameEnv        = "ARM_TEST_HPC_AD_DOMAIN_NAME"
	adCacheNetBiosNameEnv  = "ARM_TEST_HPC_AD_CACHE_NET_BIOS_NAME"
	adDomainNetBiosNameEnv = "ARM_TEST_HPC_AD_DOMAIN_NET_BIOS_NAME"
	adUsernameEnv          = "ARM_TEST_HPC_AD_USERNAME"
	adPasswordEnv          = "ARM_TEST_HPC_AD_PASSWORD"
)

func TestAccHPCCache_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache", "test")
	r := HPCCacheResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mount_addresses.#").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccHPCCache_mtu(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache", "test")
	r := HPCCacheResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.mtu(data, 1000),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mount_addresses.#").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.mtu(data, 1500),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mount_addresses.#").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.mtu(data, 1000),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mount_addresses.#").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccHPCCache_ntpServer(t *testing.T) {
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
			Config: r.ntpServer(data, "time.microsoft.com"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mount_addresses.#").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mount_addresses.#").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccHPCCache_dnsSetting(t *testing.T) {
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
			Config: r.dnsSetting(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mount_addresses.#").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mount_addresses.#").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccHPCCache_rootSquashDeprecated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache", "test")
	r := HPCCacheResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.rootSquashDeprecated(data, false),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mount_addresses.#").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.rootSquashDeprecated(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mount_addresses.#").Exists(),
			),
		},
		// Following import verification will cause diff given we simply set whatever is in cfg to state for "root_squash_enabled", since there is no
		// "cfg" during import verification, the state of the "root_squash_enabled" is always false.
		// The clarification is that since this is a deprecated property, users shouldn't import an existing resource to a new .tf file whilst using that
		// deprecated property.
		// data.ImportStep(),
		{
			Config: r.rootSquashDeprecated(data, false),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mount_addresses.#").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccHPCCache_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache", "test")
	r := HPCCacheResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccHPCCache_defaultAccessPolicy(t *testing.T) {
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

func TestAccHPCCache_updateTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache", "test")
	r := HPCCacheResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mount_addresses.#").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateTags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mount_addresses.#").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccHPCCache_directoryAD(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache", "test")
	r := HPCCacheResource{}

	requiredEnvVars := []string{adSubnetIDEnv, adPrimaryDNSEnv, adDomainNameEnv, adCacheNetBiosNameEnv, adDomainNetBiosNameEnv, adUsernameEnv, adPasswordEnv}
	for _, ev := range requiredEnvVars {
		if os.Getenv(ev) == "" {
			t.Skipf("Skip since env var %q is not set", ev)
		}
	}

	subnetId, err := networkParse.SubnetID(os.Getenv(adSubnetIDEnv))
	if err != nil {
		t.Fatal(err)
	}

	adInfo := HPCCacheDirectoryADInfo{
		SubnetID:          *subnetId,
		PrimaryDNS:        os.Getenv(adPrimaryDNSEnv),
		DomainName:        os.Getenv(adDomainNameEnv),
		CacheNetBiosName:  os.Getenv(adCacheNetBiosNameEnv),
		DomainNetBiosName: os.Getenv(adDomainNetBiosNameEnv),
		Username:          os.Getenv(adUsernameEnv),
		Password:          os.Getenv(adPasswordEnv),
	}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.adNone(data, adInfo),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.ad(data, adInfo),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("directory_active_directory.0.username", "directory_active_directory.0.password"),
		{
			Config: r.adNone(data, adInfo),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccHPCCache_directoryLDAP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache", "test")
	r := HPCCacheResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.ldapBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mount_addresses.#").Exists(),
			),
		},
		data.ImportStep("directory_ldap.0.bind"),
		{
			Config: r.ldapNone(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mount_addresses.#").Exists(),
			),
		},
		data.ImportStep("directory_ldap.0.bind"),
		{
			Config: r.ldapComplete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mount_addresses.#").Exists(),
			),
		},
		data.ImportStep("directory_ldap.0.bind"),
	})
}

func TestAccHPCCache_directoryFlatFile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hpc_cache", "test")
	r := HPCCacheResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.flatFileNone(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mount_addresses.#").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.flatFileComplete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mount_addresses.#").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.flatFileNone(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("mount_addresses.#").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func (HPCCacheResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.CacheID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.HPCCache.CachesClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving HPC Cache (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.CacheProperties != nil), nil
}

func (r HPCCacheResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache" "test" {
  name                = "acctest-HPCC-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cache_size_in_gb    = 3072
  subnet_id           = azurerm_subnet.test.id
  sku_name            = "Standard_2G"
  tags = {
    environment = "Production"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r HPCCacheResource) updateTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache" "test" {
  name                = "acctest-HPCC-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cache_size_in_gb    = 3072
  subnet_id           = azurerm_subnet.test.id
  sku_name            = "Standard_2G"
  tags = {
    environment = "Test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r HPCCacheResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache" "import" {
  name                = azurerm_hpc_cache.test.name
  resource_group_name = azurerm_hpc_cache.test.resource_group_name
  location            = azurerm_hpc_cache.test.location
  cache_size_in_gb    = azurerm_hpc_cache.test.cache_size_in_gb
  subnet_id           = azurerm_hpc_cache.test.subnet_id
  sku_name            = azurerm_hpc_cache.test.sku_name
}
`, r.basic(data))
}

func (r HPCCacheResource) mtu(data acceptance.TestData, mtu int) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache" "test" {
  name                = "acctest-HPCC-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cache_size_in_gb    = 3072
  subnet_id           = azurerm_subnet.test.id
  sku_name            = "Standard_2G"
  mtu                 = %d
}
`, r.template(data), data.RandomInteger, mtu)
}

func (r HPCCacheResource) rootSquashDeprecated(data acceptance.TestData, enable bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache" "test" {
  name                = "acctest-HPCC-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cache_size_in_gb    = 3072
  subnet_id           = azurerm_subnet.test.id
  sku_name            = "Standard_2G"
  root_squash_enabled = %t
}
`, r.template(data), data.RandomInteger, enable)
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
      scope  = "default"
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
      scope  = "default"
      access = "ro"
    }

    access_rule {
      scope                   = "network"
      access                  = "rw"
      filter                  = "10.0.0.0/24"
      suid_enabled            = true
      submount_access_enabled = true
      root_squash_enabled     = true
      anonymous_uid           = 123
      anonymous_gid           = 123
    }

    access_rule {
      scope  = "host"
      access = "no"
      filter = "10.0.0.1"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r HPCCacheResource) ntpServer(data acceptance.TestData, server string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache" "test" {
  name                = "acctest-HPCC-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cache_size_in_gb    = 3072
  subnet_id           = azurerm_subnet.test.id
  sku_name            = "Standard_2G"
  ntp_server          = %q
}
`, r.template(data), data.RandomInteger, server)
}

func (r HPCCacheResource) dnsSetting(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache" "test" {
  name                = "acctest-HPCC-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cache_size_in_gb    = 3072
  subnet_id           = azurerm_subnet.test.id
  sku_name            = "Standard_2G"
  dns {
    servers       = ["8.8.8.8"]
    search_domain = "foo.com"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r HPCCacheResource) adNone(data acceptance.TestData, info HPCCacheDirectoryADInfo) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_resource_group" "test" {
  name = "%s"
}

data "azurerm_subnet" "test" {
  resource_group_name  = "%s"
  virtual_network_name = "%s"
  name                 = "%s"
}

resource "azurerm_hpc_cache" "test" {
  name                = "acctest-HPC-%d"
  resource_group_name = data.azurerm_resource_group.test.name
  location            = data.azurerm_resource_group.test.location
  cache_size_in_gb    = 3072
  subnet_id           = data.azurerm_subnet.test.id
  sku_name            = "Standard_2G"
}
`, info.SubnetID.ResourceGroup, info.SubnetID.ResourceGroup, info.SubnetID.VirtualNetworkName, info.SubnetID.Name,
		data.RandomInteger)
}

func (r HPCCacheResource) ad(data acceptance.TestData, info HPCCacheDirectoryADInfo) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_resource_group" "test" {
  name = "%s"
}

data "azurerm_subnet" "test" {
  resource_group_name  = "%s"
  virtual_network_name = "%s"
  name                 = "%s"
}

resource "azurerm_hpc_cache" "test" {
  name                = "acctest-HPC-%d"
  resource_group_name = data.azurerm_resource_group.test.name
  location            = data.azurerm_resource_group.test.location
  cache_size_in_gb    = 3072
  subnet_id           = data.azurerm_subnet.test.id
  sku_name            = "Standard_2G"
  directory_active_directory {
    dns_primary_ip      = "%s"
    domain_name         = "%s"
    cache_netbios_name  = "%s"
    domain_netbios_name = "%s"
    username            = "%s"
    password            = "%s"
  }
}
`, info.SubnetID.ResourceGroup, info.SubnetID.ResourceGroup, info.SubnetID.VirtualNetworkName, info.SubnetID.Name,
		data.RandomInteger, info.PrimaryDNS, info.DomainName, info.CacheNetBiosName, info.DomainNetBiosName, info.Username, info.Password)
}

func (r HPCCacheResource) flatFileNone(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache" "test" {
  name                = "acctest-HPC-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cache_size_in_gb    = 3072
  subnet_id           = azurerm_subnet.test.id
  sku_name            = "Standard_2G"

  depends_on = [azurerm_linux_virtual_machine.test]
}
`, r.directoryFlatFileTemplate(data), data.RandomInteger)
}

func (r HPCCacheResource) flatFileComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache" "test" {
  name                = "acctest-HPC-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cache_size_in_gb    = 3072
  subnet_id           = azurerm_subnet.test.id
  sku_name            = "Standard_2G"
  directory_flat_file {
    group_file_uri    = "http://${azurerm_network_interface.test.private_ip_address}:8000/group"
    password_file_uri = "http://${azurerm_network_interface.test.private_ip_address}:8000/passwd"
  }

  depends_on = [azurerm_linux_virtual_machine.test]
}
`, r.directoryFlatFileTemplate(data), data.RandomInteger)
}

func (r HPCCacheResource) directoryFlatFileTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# Following script spins up a http server to host files under /etc.
locals {
  custom_data = <<CUSTOMDATA
#!/bin/bash

sudo -i
cd /etc && nohup python3 -m http.server 8000 &

CUSTOMDATA
}

resource "azurerm_linux_virtual_machine" "test" {
  name                            = "acctest-vm-%d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  size                            = "Standard_F2"
  admin_username                  = "adminuser"
  admin_password                  = "P@$$w0rd1234!"
  disable_password_authentication = false
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
    version   = "latest"
  }
  custom_data = base64encode(local.custom_data)
}
`, r.directoryTemplate(data), data.RandomInteger)
}

func (r HPCCacheResource) ldapNone(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache" "test" {
  name                = "acctest-HPC-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cache_size_in_gb    = 3072
  subnet_id           = azurerm_subnet.test.id
  sku_name            = "Standard_2G"

  depends_on = [azurerm_linux_virtual_machine.test]
}
`, r.directoryLdapTemplate(data), data.RandomInteger)
}

func (r HPCCacheResource) ldapBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache" "test" {
  name                = "acctest-HPC-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cache_size_in_gb    = 3072
  subnet_id           = azurerm_subnet.test.id
  sku_name            = "Standard_2G"
  directory_ldap {
    server  = azurerm_network_interface.test.private_ip_address
    base_dn = "dc=example,dc=com"
    bind {
      dn       = "cn=admin,dc=example,dc=com"
      password = "123"
    }
  }

  depends_on = [azurerm_linux_virtual_machine.test]
}
`, r.directoryLdapTemplate(data), data.RandomInteger)
}

func (r HPCCacheResource) ldapComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hpc_cache" "test" {
  name                = "acctest-HPC-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cache_size_in_gb    = 3072
  subnet_id           = azurerm_subnet.test.id
  sku_name            = "Standard_2G"
  directory_ldap {
    server                             = azurerm_network_interface.test.private_ip_address
    base_dn                            = "dc=example,dc=com"
    encrypted                          = true
    certificate_validation_uri         = "http://${azurerm_network_interface.test.private_ip_address}:8000/server.crt"
    download_certificate_automatically = true
    bind {
      dn       = "cn=admin,dc=example,dc=com"
      password = "123"
    }
  }

  depends_on = [azurerm_linux_virtual_machine.test]
}
`, r.directoryLdapTemplate(data), data.RandomInteger)
}

func (r HPCCacheResource) directoryLdapTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# Following script will setup a LDAP server with following settings:
# - base dn: "dc=example,dc=com"
# - admin dn: "cn=admin,dc=example,dc=com"
# - admin password: "123"
# - server cert url: http://<ip>:8000/server.crt
locals {
  custom_data = <<CUSTOMDATA
#!/bin/bash

sudo -i

hostnamectl set-hostname ldap.example.com

# Install (without specifying the root pw as we are in noninteractive mode)
DEBIAN_FRONTEND=noninteractive apt install -y slapd ldap-utils

# Update the root pw to "123"
cat << EOF > /tmp/rpw.ldif
dn: olcDatabase={1}mdb,cn=config
changetype: modify
replace: olcRootPW
olcRootPW: $(slappasswd -s 123)
EOF

ldapmodify -Q -Y EXTERNAL -H ldapi:/// -f /tmp/rpw.ldif

# Setup self signed certificate
cp /etc/ssl/certs/ca-certificates.crt /etc/ldap/sasl2
cd /etc/ldap/sasl2
openssl req -new -newkey rsa:4096 -days 365 -nodes -x509 -subj "/C=CN/ST=SH/L=SH/O=NA/CN=www.example.com" -keyout server.key -out server.crt
chown openldap. /etc/ldap/sasl2/*

cat << EOF > /tmp/cert.ldif
dn: cn=config
changetype: modify
add: olcTLSCACertificateFile
olcTLSCACertificateFile: /etc/ldap/sasl2/ca-certificates.crt
-
replace: olcTLSCertificateFile
olcTLSCertificateFile: /etc/ldap/sasl2/server.crt
-
replace: olcTLSCertificateKeyFile
olcTLSCertificateKeyFile: /etc/ldap/sasl2/server.key
EOF

ldapmodify -Q -Y EXTERNAL -H ldapi:/// -f /tmp/cert.ldif

# Host the certificate file
[[ ! -d cert ]] && mkdir /cert
cd /cert
cp /etc/ldap/sasl2/server.crt .
nohup python3 -m http.server 8000 &

CUSTOMDATA
}

resource "azurerm_linux_virtual_machine" "test" {
  name                            = "acctest-vm-%d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  size                            = "Standard_F2"
  admin_username                  = "adminuser"
  admin_password                  = "P@$$w0rd1234!"
  disable_password_authentication = false
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
    version   = "latest"
  }
  custom_data = base64encode(local.custom_data)
}
`, r.directoryTemplate(data), data.RandomInteger)
}

func (HPCCacheResource) directoryTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VN-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (HPCCacheResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VN-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsub-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
