package azurerm

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAdApplication_simple(t *testing.T) {
	resourceName := "azurerm_ad_application.test"
	id := uuid.New().String()
	config := testAccAzureRMAdApplication_simple(id)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAdApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdApplicationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", id),
					resource.TestCheckResourceAttr(resourceName, "homepage", fmt.Sprintf("http://%s", id)),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName, "object_id"),
				),
			},
		},
	})
}

func TestAccAzureRMAdApplication_advanced(t *testing.T) {
	resourceName := "azurerm_ad_application.test"
	id := uuid.New().String()
	config := testAccAzureRMAdApplication_advanced(id)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAdApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdApplicationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", id),
					resource.TestCheckResourceAttr(resourceName, "homepage", fmt.Sprintf("http://homepage-%s", id)),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName, "object_id"),
				),
			},
		},
	})
}

func TestAccAzureRMAdApplication_keyCredential(t *testing.T) {
	resourceName := "azurerm_ad_application.test"
	id := uuid.New().String()
	keyId := uuid.New().String()
	config := testAccAzureRMAdApplication_keyCredential(id, keyId)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProvidersWithTLS,
		CheckDestroy: testCheckAzureRMAdApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdApplicationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", id),
					resource.TestCheckResourceAttr(resourceName, "homepage", fmt.Sprintf("http://%s", id)),
					resource.TestCheckResourceAttr(resourceName, "key_credential.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName, "object_id"),
				),
			},
		},
	})
}

func TestAccAzureRMAdApplication_passwordCredential(t *testing.T) {
	resourceName := "azurerm_ad_application.test"
	id := uuid.New().String()
	keyId := uuid.New().String()
	timeStart := string(time.Now().UTC().Format(time.RFC3339))
	timeEnd := string(time.Now().UTC().Add(time.Duration(1) * time.Hour).Format(time.RFC3339))
	config := testAccAzureRMAdApplication_passwordCredential(id, keyId, timeStart, timeEnd)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAdApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdApplicationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", id),
					resource.TestCheckResourceAttr(resourceName, "homepage", fmt.Sprintf("http://%s", id)),
					resource.TestCheckResourceAttr(resourceName, "password_credential.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName, "object_id"),
				),
			},
		},
	})
}

func testCheckAzureRMAdApplicationExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %q", name)
		}

		objectId := rs.Primary.Attributes["object_id"]

		client := testAccProvider.Meta().(*ArmClient).applicationsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, objectId)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: AD Application %q does not exist", objectId)
			}
			return fmt.Errorf("Bad: Get on applicationsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMAdApplicationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_ad_application" {
			continue
		}

		objectId := rs.Primary.Attributes["object_id"]

		client := testAccProvider.Meta().(*ArmClient).applicationsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, objectId)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("AD Application still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMAdApplication_simple(id string) string {
	return fmt.Sprintf(`
resource "azurerm_ad_application" "test" {
  display_name = "%s"
}
`, id)
}

func testAccAzureRMAdApplication_advanced(id string) string {
	return fmt.Sprintf(`
resource "azurerm_ad_application" "test" {
  display_name = "%s"
  homepage = "http://homepage-%s"
  identifier_uris = ["http://uri-%s"]
  reply_urls = ["http://replyurl-%s"]
  available_to_other_tenants = false
  oauth2_allow_implicit_flow = true
}
`, id, id, id, id)
}

func testAccAzureRMAdApplication_keyCredential(id string, keyId string) string {
	return fmt.Sprintf(`
resource "tls_private_key" "test" {
  algorithm   = "ECDSA"
  ecdsa_curve = "P384"
}

resource "tls_self_signed_cert" "test" {
  key_algorithm   = "${tls_private_key.test.algorithm}"
  private_key_pem = "${tls_private_key.test.private_key_pem}"

  subject {
    common_name  = "example.com"
    organization = "ACME Examples, Inc"
  }

  validity_period_hours = 12

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
    "cert_signing",
  ]
}

resource "azurerm_ad_application" "test" {
  display_name = "%s"

  key_credential {
    key_id = "%s"
    type = "AsymmetricX509Cert"
    usage = "Verify"
    value = "${replace(tls_self_signed_cert.test.cert_pem, "/(-{5}.+?-{5})|(\\n)/", "")}"
  }
}
`, id, keyId)
}

func testAccAzureRMAdApplication_passwordCredential(id string, keyId string, timeStart string, timeEnd string) string {
	return fmt.Sprintf(`
resource "azurerm_ad_application" "test" {
  display_name = "%s"

  password_credential {
    key_id = "%s"
    value = "test"
    start_date = "%s"
    end_date = "%s"
  }
}
`, id, keyId, timeStart, timeEnd)
}
