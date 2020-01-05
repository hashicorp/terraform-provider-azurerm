package tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault"
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
		_, errors := keyvault.ValidateKeyVaultName(tc.Input, "")

		hasError := len(errors) > 0
		if tc.ExpectError && !hasError {
			t.Fatalf("Expected the Key Vault Name to trigger a validation error for '%s'", tc.Input)
		}
	}
}

func TestAccAzureRMKeyVault_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")
	config := testAccAzureRMKeyVault_basic(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "premium"),
				),
			},
			data.ImportStep(),
		},
	})
}

// Remove in 2.0
func TestAccAzureRMKeyVault_basicNotDefined(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")
	config := testAccAzureRMKeyVault_basicNotDefined(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")
	config := testAccAzureRMKeyVault_basicClassic(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "premium"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKeyVault_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVault_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMKeyVault_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_key_vault"),
			},
		},
	})
}

func TestAccAzureRMKeyVault_networkAcls(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVault_networkAcls(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMKeyVault_networkAclsUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKeyVault_networkAclsAllowed(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVault_networkAclsAllowed(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKeyVault_accessPolicyUpperLimit(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVault_accessPolicyUpperLimit(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultExists(data.ResourceName),
					testCheckAzureRMKeyVaultDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMKeyVault_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVault_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultExists(data.ResourceName),
					testCheckAzureRMKeyVaultDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMKeyVault_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVault_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "access_policy.0.application_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKeyVault_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")
	preConfig := testAccAzureRMKeyVault_basic(data)
	postConfig := testAccAzureRMKeyVault_update(data)
	noAccessPolicyConfig := testAccAzureRMKeyVault_noAccessPolicyBlocks(data)
	forceZeroAccessPolicyConfig := testAccAzureRMKeyVault_accessPolicyExplicitZero(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "access_policy.0.key_permissions.0", "create"),
					resource.TestCheckResourceAttr(data.ResourceName, "access_policy.0.secret_permissions.0", "set"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "access_policy.0.key_permissions.0", "get"),
					resource.TestCheckResourceAttr(data.ResourceName, "access_policy.0.secret_permissions.0", "get"),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled_for_deployment", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled_for_disk_encryption", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled_for_template_deployment", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "Staging"),
				),
			},
			{
				Config: noAccessPolicyConfig,
				Check: resource.ComposeTestCheckFunc(
					// There are no access_policy blocks in this configuration
					// at all, which means to ignore any existing policies and
					// so the one created in previous steps is still present.
					resource.TestCheckResourceAttr(data.ResourceName, "access_policy.#", "1"),
				),
			},
			{
				Config: forceZeroAccessPolicyConfig,
				Check: resource.ComposeTestCheckFunc(
					// This config explicitly sets access_policy = [], which
					// means to delete any existing policies.
					resource.TestCheckResourceAttr(data.ResourceName, "access_policy.#", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMKeyVault_justCert(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVault_justCert(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "access_policy.0.certificate_permissions.0", "get"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMKeyVaultDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.VaultsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

		client := acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.VaultsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

		client := acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.VaultsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testAccAzureRMKeyVault_basic(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMKeyVault_basicNotDefined(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMKeyVault_basicClassic(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMKeyVault_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMKeyVault_basic(data)
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

func testAccAzureRMKeyVault_networkAclsTemplate(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMKeyVault_networkAcls(data acceptance.TestData) string {
	template := testAccAzureRMKeyVault_networkAclsTemplate(data)
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
`, template, data.RandomInteger)
}

func testAccAzureRMKeyVault_networkAclsUpdated(data acceptance.TestData) string {
	template := testAccAzureRMKeyVault_networkAclsTemplate(data)
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
    ip_rules                   = ["123.0.0.102/32"]
    virtual_network_subnet_ids = ["${azurerm_subnet.test_a.id}"]
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMKeyVault_networkAclsAllowed(data acceptance.TestData) string {
	template := testAccAzureRMKeyVault_networkAclsTemplate(data)
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
    default_action = "Allow"
    bypass         = "AzureServices"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMKeyVault_update(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMKeyVault_noAccessPolicyBlocks(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMKeyVault_accessPolicyExplicitZero(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMKeyVault_complete(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMKeyVault_justCert(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMKeyVault_accessPolicyUpperLimit(data acceptance.TestData) string {
	var storageAccountConfigs string
	var accessPoliciesConfigs string

	for i := 1; i <= 20; i++ {
		storageAccountConfigs += testAccAzureRMKeyVault_generateStorageAccountConfigs(i, data.RandomString)
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

`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, accessPoliciesConfigs, storageAccountConfigs)
}

func testAccAzureRMKeyVault_generateStorageAccountConfigs(accountNum int, rs string) string {
	return fmt.Sprintf(`
resource "azurerm_storage_account" "test%d" {
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
  object_id = "${azurerm_storage_account.test%d.identity.0.principal_id}"

  key_permissions    = ["get", "create", "delete", "list", "restore", "recover", "unwrapkey", "wrapkey", "purge", "encrypt", "decrypt", "sign", "verify"]
  secret_permissions = ["get"]
}
`, accountNum)
}
