package azurerm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMKeyVaultAccessPolicy_basic(t *testing.T) {
	resourceName := "azurerm_key_vault_access_policy.test"
	rs := acctest.RandString(6)
	config := testAccAzureRMKeyVaultAccessPolicy_basic(rs, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultAccessPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "key_permissions.0", "get"),
					resource.TestCheckResourceAttr(resourceName, "secret_permissions.0", "get"),
					resource.TestCheckResourceAttr(resourceName, "secret_permissions.1", "set"),
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

func TestAccAzureRMKeyVaultAccessPolicy_basicClassic(t *testing.T) {
	resourceName := "azurerm_key_vault_access_policy.test"
	rs := acctest.RandString(6)
	config := testAccAzureRMKeyVaultAccessPolicy_basicClassic(rs, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultAccessPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "key_permissions.0", "get"),
					resource.TestCheckResourceAttr(resourceName, "secret_permissions.0", "get"),
					resource.TestCheckResourceAttr(resourceName, "secret_permissions.1", "set"),
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

func TestAccAzureRMKeyVaultAccessPolicy_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_key_vault_access_policy.test"
	rs := acctest.RandString(6)
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKeyVaultAccessPolicy_basic(rs, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultAccessPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "key_permissions.0", "get"),
					resource.TestCheckResourceAttr(resourceName, "secret_permissions.0", "get"),
					resource.TestCheckResourceAttr(resourceName, "secret_permissions.1", "set"),
				),
			},
			{
				Config:      testAccAzureRMKeyVaultAccessPolicy_requiresImport(rs, location),
				ExpectError: acceptance.RequiresImportError("azurerm_key_vault_access_policy"),
			},
		},
	})
}

func TestAccAzureRMKeyVaultAccessPolicy_multiple(t *testing.T) {
	resourceName1 := "azurerm_key_vault_access_policy.test_with_application_id"
	resourceName2 := "azurerm_key_vault_access_policy.test_no_application_id"
	rs := acctest.RandString(6)
	config := testAccAzureRMKeyVaultAccessPolicy_multiple(rs, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultAccessPolicyExists(resourceName1),
					resource.TestCheckResourceAttr(resourceName1, "key_permissions.0", "create"),
					resource.TestCheckResourceAttr(resourceName1, "key_permissions.1", "get"),
					resource.TestCheckResourceAttr(resourceName1, "secret_permissions.0", "get"),
					resource.TestCheckResourceAttr(resourceName1, "secret_permissions.1", "delete"),
					resource.TestCheckResourceAttr(resourceName1, "certificate_permissions.0", "create"),
					resource.TestCheckResourceAttr(resourceName1, "certificate_permissions.1", "delete"),

					testCheckAzureRMKeyVaultAccessPolicyExists(resourceName2),
					resource.TestCheckResourceAttr(resourceName2, "key_permissions.0", "list"),
					resource.TestCheckResourceAttr(resourceName2, "key_permissions.1", "encrypt"),
					resource.TestCheckResourceAttr(resourceName2, "secret_permissions.0", "list"),
					resource.TestCheckResourceAttr(resourceName2, "secret_permissions.1", "delete"),
					resource.TestCheckResourceAttr(resourceName2, "certificate_permissions.0", "list"),
					resource.TestCheckResourceAttr(resourceName2, "certificate_permissions.1", "delete"),
				),
			},
			{
				ResourceName:      resourceName1,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      resourceName2,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMKeyVaultAccessPolicy_update(t *testing.T) {
	rs := acctest.RandString(6)
	resourceName := "azurerm_key_vault_access_policy.test"
	preConfig := testAccAzureRMKeyVaultAccessPolicy_basic(rs, acceptance.Location())
	postConfig := testAccAzureRMKeyVaultAccessPolicy_update(rs, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultAccessPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "key_permissions.0", "get"),
					resource.TestCheckResourceAttr(resourceName, "secret_permissions.0", "get"),
					resource.TestCheckResourceAttr(resourceName, "secret_permissions.1", "set"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKeyVaultAccessPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "key_permissions.0", "list"),
					resource.TestCheckResourceAttr(resourceName, "key_permissions.1", "encrypt"),
				),
			},
		},
	})
}

func TestAccAzureRMKeyVaultAccessPolicy_nonExistentVault(t *testing.T) {
	rs := acctest.RandString(6)
	config := testAccAzureRMKeyVaultAccessPolicy_nonExistentVault(rs, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKeyVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config:             config,
				ExpectNonEmptyPlan: true,
				ExpectError:        regexp.MustCompile(`Error retrieving Key Vault`),
			},
		},
	})
}

func testCheckAzureRMKeyVaultAccessPolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).KeyVault.VaultsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := azure.ParseAzureResourceID(rs.Primary.ID)

		if err != nil {
			return err
		}
		resGroup := id.ResourceGroup
		vaultName := id.Path["vaults"]

		objectId := rs.Primary.Attributes["object_id"]
		applicationId := rs.Primary.Attributes["application_id"]

		resp, err := client.Get(ctx, resGroup, vaultName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Key Vault %q (resource group: %q) does not exist", vaultName, resGroup)
			}

			return fmt.Errorf("Bad: Get on keyVaultClient: %+v", err)
		}

		policy, err := findKeyVaultAccessPolicy(resp.Properties.AccessPolicies, objectId, applicationId)
		if err != nil {
			return fmt.Errorf("Error finding Key Vault Access Policy %q : %+v", vaultName, err)
		}
		if policy == nil {
			return fmt.Errorf("Bad: Key Vault Policy %q (resource group: %q, object_id: %s) does not exist", vaultName, resGroup, objectId)
		}

		return nil
	}
}

func testAccAzureRMKeyVaultAccessPolicy_basic(rString string, location string) string {
	template := testAccAzureRMKeyVaultAccessPolicy_template(rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = "${azurerm_key_vault.test.id}"

  key_permissions = [
    "get",
  ]

  secret_permissions = [
    "get",
    "set",
  ]

  tenant_id = "${data.azurerm_client_config.current.tenant_id}"
  object_id = "${data.azurerm_client_config.current.service_principal_object_id}"
}
`, template)
}

func testAccAzureRMKeyVaultAccessPolicy_basicClassic(rString string, location string) string {
	template := testAccAzureRMKeyVaultAccessPolicy_template(rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_access_policy" "test" {
  vault_name          = "${azurerm_key_vault.test.name}"
  resource_group_name = "${azurerm_key_vault.test.resource_group_name}"

  key_permissions = [
    "get",
  ]

  secret_permissions = [
    "get",
    "set",
  ]

  tenant_id = "${data.azurerm_client_config.current.tenant_id}"
  object_id = "${data.azurerm_client_config.current.service_principal_object_id}"
}
`, template)
}

func testAccAzureRMKeyVaultAccessPolicy_requiresImport(rString string, location string) string {
	template := testAccAzureRMKeyVaultAccessPolicy_basic(rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_access_policy" "import" {
  key_vault_id = "${azurerm_key_vault.test.id}"
  tenant_id    = "${azurerm_key_vault_access_policy.test.tenant_id}"
  object_id    = "${azurerm_key_vault_access_policy.test.object_id}"

  key_permissions = [
    "get",
  ]

  secret_permissions = [
    "get",
    "set",
  ]
}
`, template)
}

func testAccAzureRMKeyVaultAccessPolicy_multiple(rString string, location string) string {
	template := testAccAzureRMKeyVaultAccessPolicy_template(rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_access_policy" "test_with_application_id" {
  key_vault_id = "${azurerm_key_vault.test.id}"

  key_permissions = [
    "create",
    "get",
  ]

  secret_permissions = [
    "get",
    "delete",
  ]

  certificate_permissions = [
    "create",
    "delete",
  ]

  application_id = "${data.azurerm_client_config.current.service_principal_application_id}"
  tenant_id      = "${data.azurerm_client_config.current.tenant_id}"
  object_id      = "${data.azurerm_client_config.current.service_principal_object_id}"
}

resource "azurerm_key_vault_access_policy" "test_no_application_id" {
  key_vault_id = "${azurerm_key_vault.test.id}"

  key_permissions = [
    "list",
    "encrypt",
  ]

  secret_permissions = [
    "list",
    "delete",
  ]

  certificate_permissions = [
    "list",
    "delete",
  ]

  storage_permissions = [
    "backup",
    "delete",
    "deletesas",
    "get",
    "getsas",
    "list",
    "listsas",
    "purge",
    "recover",
    "regeneratekey",
    "restore",
    "set",
    "setsas",
    "update",
  ]

  tenant_id = "${data.azurerm_client_config.current.tenant_id}"
  object_id = "${data.azurerm_client_config.current.service_principal_object_id}"
}
`, template)
}

func testAccAzureRMKeyVaultAccessPolicy_update(rString string, location string) string {
	template := testAccAzureRMKeyVaultAccessPolicy_template(rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_access_policy" "test" {
  vault_name          = "${azurerm_key_vault.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  key_permissions = [
    "list",
    "encrypt",
  ]

  secret_permissions = []

  tenant_id = "${data.azurerm_client_config.current.tenant_id}"
  object_id = "${data.azurerm_client_config.current.service_principal_object_id}"
}
`, template)
}

func testAccAzureRMKeyVaultAccessPolicy_template(rString string, location string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%s"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku_name = "premium"

  tags = {
    environment = "Production"
  }
}
`, rString, location, rString)
}

func testAccAzureRMKeyVaultAccessPolicy_nonExistentVault(rString string, location string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%s"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku_name = "standard"

  tags = {
    environment = "Production"
  }
}

resource "azurerm_key_vault_access_policy" "test" {
  # Must appear to be URL, but not actually exist - appending a string works
  key_vault_id = "${azurerm_key_vault.test.id}NOPE"

  tenant_id = "${data.azurerm_client_config.current.tenant_id}"
  object_id = "${data.azurerm_client_config.current.service_principal_object_id}"

  key_permissions = [
    "get",
  ]

  secret_permissions = [
    "get",
  ]
}
`, rString, location, rString)
}
