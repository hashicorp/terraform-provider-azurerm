package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMKeyVault_name(t *testing.T) {
	cases := []struct {
		Input       string
		ExpectError bool
	}{
		{
			Input:       "",
			ExpectError: true,
		},
		{
			Input:       "hi",
			ExpectError: true,
		},
		{
			Input:       "hello",
			ExpectError: false,
		},
		{
			Input:       "hello-world",
			ExpectError: false,
		},
		{
			Input:       "hello-world-21",
			ExpectError: false,
		},
		{
			Input:       "hello_world_21",
			ExpectError: true,
		},
		{
			Input:       "Hello-World",
			ExpectError: false,
		},
		{
			Input:       "20202020",
			ExpectError: false,
		},
		{
			Input:       "ABC123!@£",
			ExpectError: true,
		},
		{
			Input:       "abcdefghijklmnopqrstuvwxyz",
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		_, errors := validateKeyVaultName(tc.Input, "")

		hasError := len(errors) > 0

		if tc.ExpectError && !hasError {
			t.Fatalf("Expected the Key Vault Name to trigger a validation error for '%s'", tc.Input)
		}
	}
}

func TestAccAzureRMKeyVault_basic(t *testing.T) {
	resourceName := "azurerm_key_vault.test"
	ri := acctest.RandInt()
	config := testAccAzureRMKeyVault_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "network_acls.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMKeyVault_networkAcls(t *testing.T) {
	resourceName := "azurerm_key_vault.test"
	ri := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVault_networkAcls(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "network_acls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_acls.0.bypass", "None"),
					resource.TestCheckResourceAttr(resourceName, "network_acls.0.default_action", "Deny"),
					resource.TestCheckResourceAttr(resourceName, "network_acls.0.ip_rules.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "network_acls.0.virtual_network_subnet_ids.#", "1"),
				),
			},
			{
				Config: testAccAzureRMKeyVault_networkAclsUpdated(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "network_acls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_acls.0.bypass", "AzureServices"),
					resource.TestCheckResourceAttr(resourceName, "network_acls.0.default_action", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "network_acls.0.ip_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_acls.0.virtual_network_subnet_ids.#", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMKeyVault_disappears(t *testing.T) {
	resourceName := "azurerm_key_vault.test"
	ri := acctest.RandInt()
	config := testAccAzureRMKeyVault_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultExists(resourceName),
					testCheckAzureRMKeyVaultDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMKeyVault_complete(t *testing.T) {
	resourceName := "azurerm_key_vault.test"
	ri := acctest.RandInt()
	config := testAccAzureRMKeyVault_complete(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "access_policy.0.application_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMKeyVault_update(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := "azurerm_key_vault.test"
	preConfig := testAccAzureRMKeyVault_basic(ri, testLocation())
	postConfig := testAccAzureRMKeyVault_update(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "access_policy.0.key_permissions.0", "create"),
					resource.TestCheckResourceAttr(resourceName, "access_policy.0.secret_permissions.0", "set"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "access_policy.0.key_permissions.0", "get"),
					resource.TestCheckResourceAttr(resourceName, "access_policy.0.secret_permissions.0", "get"),
					resource.TestCheckResourceAttr(resourceName, "enabled_for_deployment", "true"),
					resource.TestCheckResourceAttr(resourceName, "enabled_for_disk_encryption", "true"),
					resource.TestCheckResourceAttr(resourceName, "enabled_for_template_deployment", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "Staging"),
				),
			},
		},
	})
}

func testCheckAzureRMKeyVaultDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).keyVaultClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_key_vault" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		return fmt.Errorf("Key Vault still exists:\n%#v", resp.Properties)
	}

	return nil
}

func testCheckAzureRMKeyVaultExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		vaultName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for vault: %s", vaultName)
		}

		client := testAccProvider.Meta().(*ArmClient).keyVaultClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, vaultName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Vault %q (resource group: %q) does not exist", vaultName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on keyVaultClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMKeyVaultDisappears(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		vaultName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for vault: %s", vaultName)
		}

		client := testAccProvider.Meta().(*ArmClient).keyVaultClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Delete(ctx, resourceGroup, vaultName)
		if err != nil {
			if response.WasNotFound(resp.Response) {
				return nil
			}

			return fmt.Errorf("Bad: Delete on keyVaultClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMKeyVault_basic(rInt int, location string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "vault%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku {
    name = "premium"
  }

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.client_id}"

    key_permissions = [
      "create",
    ]

    secret_permissions = [
      "set",
    ]
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMKeyVault_networkAclsTemplate(rInt int, location string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
  service_endpoints    = ["Microsoft.KeyVault",]
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMKeyVault_networkAcls(rInt int, location string) string {
	template := testAccAzureRMKeyVault_networkAclsTemplate(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault" "test" {
  name                = "vault%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku {
    name = "premium"
  }

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.client_id}"

    key_permissions = [
      "create",
    ]

    secret_permissions = [
      "set",
    ]
  }

  network_acls {
    default_action             = "Deny"
    bypass                     = "None"
    virtual_network_subnet_ids = ["${azurerm_subnet.test.id}"]
  }
}
`, template, rInt)
}

func testAccAzureRMKeyVault_networkAclsUpdated(rInt int, location string) string {
	template := testAccAzureRMKeyVault_networkAclsTemplate(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault" "test" {
  name                = "vault%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku {
    name = "premium"
  }

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.client_id}"

    key_permissions = [
      "create",
    ]

    secret_permissions = [
      "set",
    ]
  }

  network_acls {
    default_action             = "Allow"
    bypass                     = "AzureServices"
    ip_rules                   = ["10.0.0.102/32"]
    virtual_network_subnet_ids = ["${azurerm_subnet.test.id}"]
  }
}
`, template, rInt)
}

func testAccAzureRMKeyVault_update(rInt int, location string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "vault%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku {
    name = "premium"
  }

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.client_id}"

    key_permissions = [
      "get",
    ]

    secret_permissions = [
      "get",
    ]
  }

  enabled_for_deployment          = true
  enabled_for_disk_encryption     = true
  enabled_for_template_deployment = true

  tags {
    environment = "Staging"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMKeyVault_complete(rInt int, location string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "vault%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku {
    name = "premium"
  }

  access_policy {
    tenant_id      = "${data.azurerm_client_config.current.tenant_id}"
    object_id      = "${data.azurerm_client_config.current.client_id}"
    application_id = "${data.azurerm_client_config.current.service_principal_application_id}"

    certificate_permissions = [
      "get",
    ]

    key_permissions = [
      "get",
    ]

    secret_permissions = [
      "get",
    ]
  }

  tags {
    environment = "Production"
  }
}
`, rInt, location, rInt)
}
