package azurerm

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMActiveDirectoryServicePrincipalPassword_basic(t *testing.T) {
	resourceName := "azurerm_azuread_service_principal_password.test"
	applicationId, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	value, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}

	config := testAccAzureRMActiveDirectoryServicePrincipalPassword_basic(applicationId, value)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryServicePrincipalDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					// can't assert on Value since it's not returned
					testCheckAzureRMActiveDirectoryServicePrincipalPasswordExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "start_date"),
					resource.TestCheckResourceAttrSet(resourceName, "key_id"),
					resource.TestCheckResourceAttr(resourceName, "end_date", "2020-01-01T01:02:03Z"),
				),
			},
		},
	})
}

func TestAccAzureRMActiveDirectoryServicePrincipalPassword_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_azuread_service_principal_password.test"
	applicationId, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	value, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryServicePrincipalDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMActiveDirectoryServicePrincipalPassword_basic(applicationId, value),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryServicePrincipalPasswordExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMActiveDirectoryServicePrincipalPassword_requiresImport(applicationId, value),
				ExpectError: acceptance.RequiresImportError("azurerm_azuread_service_principal_password"),
			},
		},
	})
}

func TestAccAzureRMActiveDirectoryServicePrincipalPassword_customKeyId(t *testing.T) {
	resourceName := "azurerm_azuread_service_principal_password.test"
	applicationId, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	keyId, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	value, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	config := testAccAzureRMActiveDirectoryServicePrincipalPassword_customKeyId(applicationId, keyId, value)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryServicePrincipalDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					// can't assert on Value since it's not returned
					testCheckAzureRMActiveDirectoryServicePrincipalPasswordExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "start_date"),
					resource.TestCheckResourceAttr(resourceName, "key_id", keyId),
					resource.TestCheckResourceAttr(resourceName, "end_date", "2020-01-01T01:02:03Z"),
				),
			},
		},
	})
}

func testCheckAzureRMActiveDirectoryServicePrincipalPasswordExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Graph.ServicePrincipalsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		id := strings.Split(rs.Primary.ID, "/")
		objectId := id[0]
		keyId := id[1]
		resp, err := client.Get(ctx, objectId)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Azure AD Service Principal %q does not exist", objectId)
			}
			return fmt.Errorf("Bad: Get on Azure AD servicePrincipalsClient: %+v", err)
		}

		credentials, err := client.ListPasswordCredentials(ctx, objectId)
		if err != nil {
			return fmt.Errorf("Error Listing Password Credentials for Service Principal %q: %+v", objectId, err)
		}

		for _, credential := range *credentials.Value {
			if credential.KeyID == nil {
				continue
			}

			if *credential.KeyID == keyId {
				return nil
			}
		}

		return fmt.Errorf("Password Credential %q was not found in Service Principal %q", keyId, objectId)
	}
}

func testAccAzureRMActiveDirectoryServicePrincipalPassword_basic(applicationId, value string) string {
	return fmt.Sprintf(`
resource "azurerm_azuread_application" "test" {
  name = "acctestspa%s"
}

resource "azurerm_azuread_service_principal" "test" {
  application_id = "${azurerm_azuread_application.test.application_id}"
}

resource "azurerm_azuread_service_principal_password" "test" {
  service_principal_id = "${azurerm_azuread_service_principal.test.id}"
  value                = "%s"
  end_date             = "2020-01-01T01:02:03Z"
}
`, applicationId, value)
}

func testAccAzureRMActiveDirectoryServicePrincipalPassword_requiresImport(applicationId, value string) string {
	template := testAccAzureRMActiveDirectoryServicePrincipalPassword_basic(applicationId, value)
	return fmt.Sprintf(`
%s

resource "azurerm_azuread_service_principal_password" "import" {
  service_principal_id = "${azurerm_azuread_service_principal_password.test.service_principal_id}"
  value                = "${azurerm_azuread_service_principal_password.test.value}"
  end_date             = "${azurerm_azuread_service_principal_password.test.end_date}"
}
`, template)
}

func testAccAzureRMActiveDirectoryServicePrincipalPassword_customKeyId(applicationId, keyId, value string) string {
	return fmt.Sprintf(`
resource "azurerm_azuread_application" "test" {
  name = "acctestspa%s"
}

resource "azurerm_azuread_service_principal" "test" {
  application_id = "${azurerm_azuread_application.test.application_id}"
}

resource "azurerm_azuread_service_principal_password" "test" {
  service_principal_id = "${azurerm_azuread_service_principal.test.id}"
  key_id               = "%s"
  value                = "%s"
  end_date             = "2020-01-01T01:02:03Z"
}
`, applicationId, keyId, value)
}
