package azurerm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
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
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMKeyVault_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "network_acls.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "sku_name", "premium"),
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

// Remove in 2.0
func TestAccAzureRMKeyVault_basicNotDefined(t *testing.T) {
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMKeyVault_basicNotDefined(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("either 'sku_name' or 'sku' must be defined in the configuration file"),
			},
		},
	})
}

// Remove in 2.0
func TestAccAzureRMKeyVault_basicClassic(t *testing.T) {
	resourceName := "azurerm_key_vault.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMKeyVault_basicClassic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "network_acls.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "sku.0.name", "premium"),
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

func TestAccAzureRMKeyVault_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_key_vault.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVault_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "network_acls.#", "0"),
				),
			},
			{
				Config:      testAccAzureRMKeyVault_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_key_vault"),
			},
		},
	})
}

func TestAccAzureRMKeyVault_networkAcls(t *testing.T) {
	resourceName := "azurerm_key_vault.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
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
					resource.TestCheckResourceAttr(resourceName, "network_acls.0.virtual_network_subnet_ids.#", "2"),
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

func TestAccAzureRMKeyVault_accessPolicyUpperLimit(t *testing.T) {
	resourceName := "azurerm_key_vault.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(10)
	config := testAccAzureRMKeyVault_accessPolicyUpperLimit(ri, testLocation(), rs)

	resource.ParallelTest(t, resource.TestCase{
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

func TestAccAzureRMKeyVault_disappears(t *testing.T) {
	resourceName := "azurerm_key_vault.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMKeyVault_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
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
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMKeyVault_complete(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
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
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_key_vault.test"
	preConfig := testAccAzureRMKeyVault_basic(ri, testLocation())
	postConfig := testAccAzureRMKeyVault_update(ri, testLocation())
	noAccessPolicyConfig := testAccAzureRMKeyVault_noAccessPolicyBlocks(ri, testLocation())
	forceZeroAccessPolicyConfig := testAccAzureRMKeyVault_accessPolicyExplicitZero(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
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
			{
				Config: noAccessPolicyConfig,
				Check: resource.ComposeTestCheckFunc(
					// There are no access_policy blocks in this configuration
					// at all, which means to ignore any existing policies and
					// so the one created in previous steps is still present.
					resource.TestCheckResourceAttr(resourceName, "access_policy.#", "1"),
				),
			},
			{
				Config: forceZeroAccessPolicyConfig,
				Check: resource.ComposeTestCheckFunc(
					// This config explicitly sets access_policy = [], which
					// means to delete any existing policies.
					resource.TestCheckResourceAttr(resourceName, "access_policy.#", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMKeyVault_justCert(t *testing.T) {
	resourceName := "azurerm_key_vault.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVault_justCert(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "access_policy.0.certificate_permissions.0", "get"),
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

func testCheckAzureRMKeyVaultDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).KeyVault.VaultsClient
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

func testCheckAzureRMKeyVaultExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		vaultName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for vault: %s", vaultName)
		}

		client := testAccProvider.Meta().(*ArmClient).KeyVault.VaultsClient
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

func testCheckAzureRMKeyVaultDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		vaultName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for vault: %s", vaultName)
		}

		client := testAccProvider.Meta().(*ArmClient).KeyVault.VaultsClient
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

  sku_name = "premium"

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

func testAccAzureRMKeyVault_basicNotDefined(rInt int, location string) string {
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

func testAccAzureRMKeyVault_basicClassic(rInt int, location string) string {
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

func testAccAzureRMKeyVault_requiresImport(rInt int, location string) string {
	template := testAccAzureRMKeyVault_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault" "import" {
  name                = "${azurerm_key_vault.test.name}"
  location            = "${azurerm_key_vault.test.location}"
  resource_group_name = "${azurerm_key_vault.test.resource_group_name}"
  tenant_id           = "${azurerm_key_vault.test.tenant_id}"

  sku_name = "premium"

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
`, template)
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

resource "azurerm_subnet" "test_a" {
  name                 = "acctestsubneta%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
  service_endpoints    = ["Microsoft.KeyVault"]
}

resource "azurerm_subnet" "test_b" {
  name                 = "acctestsubnetb%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.4.0/24"
  service_endpoints    = ["Microsoft.KeyVault"]
}
`, rInt, location, rInt, rInt, rInt)
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

  sku_name = "premium"

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
    virtual_network_subnet_ids = ["${azurerm_subnet.test_a.id}", "${azurerm_subnet.test_b.id}"]
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

  sku_name = "premium"

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
    virtual_network_subnet_ids = ["${azurerm_subnet.test_a.id}"]
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

  sku_name = "premium"

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

  tags = {
    environment = "Staging"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMKeyVault_noAccessPolicyBlocks(rInt int, location string) string {
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

  enabled_for_deployment          = true
  enabled_for_disk_encryption     = true
  enabled_for_template_deployment = true

  tags = {
    environment = "Staging"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMKeyVault_accessPolicyExplicitZero(rInt int, location string) string {
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

  access_policy = []

  enabled_for_deployment          = true
  enabled_for_disk_encryption     = true
  enabled_for_template_deployment = true

  tags = {
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

  sku_name = "premium"

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

  tags = {
    environment = "Production"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMKeyVault_justCert(rInt int, location string) string {
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

  sku_name = "premium"

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.client_id}"

    certificate_permissions = [
      "get",
    ]
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMKeyVault_accessPolicyUpperLimit(rInt int, location string, rs string) string {
	var storageAccountConfigs string
	var accessPoliciesConfigs string

	for i := 1; i <= 20; i++ {
		storageAccountConfigs += testAccAzureRMKeyVault_generateStorageAccountConfigs(i, rs)
		accessPoliciesConfigs += testAccAzureRMKeyVault_generateAccessPolicyConfigs(i)
	}

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

  sku_name = "premium"
  %s
}

%s

`, rInt, location, rInt, accessPoliciesConfigs, storageAccountConfigs)
}

func testAccAzureRMKeyVault_generateStorageAccountConfigs(accountNum int, rs string) string {
	return fmt.Sprintf(`
resource "azurerm_storage_account" "testsa%d" {
  name                     = "testsa%s%d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"

  identity {
    type = "SystemAssigned"
  }

  tags = {
    environment = "testing"
  }
}
`, accountNum, rs, accountNum)
}

func testAccAzureRMKeyVault_generateAccessPolicyConfigs(accountNum int) string {
	return fmt.Sprintf(`
access_policy {
  tenant_id = "${data.azurerm_client_config.current.tenant_id}"
  object_id = "${azurerm_storage_account.testsa%d.identity.0.principal_id}"

  key_permissions    = ["get", "create", "delete", "list", "restore", "recover", "unwrapkey", "wrapkey", "purge", "encrypt", "decrypt", "sign", "verify"]
  secret_permissions = ["get"]
}
`, accountNum)
}
