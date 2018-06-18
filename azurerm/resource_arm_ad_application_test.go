package azurerm

import (
	"fmt"
	"regexp"
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

func TestAccAzureRMAdApplication_updateAdvanced(t *testing.T) {
	resourceName := "azurerm_ad_application.test"
	id := uuid.New().String()
	config := testAccAzureRMAdApplication_simple(id)

	updatedId := uuid.New().String()
	updatedConfig := testAccAzureRMAdApplication_advanced(updatedId)

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
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdApplicationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", updatedId),
					resource.TestCheckResourceAttr(resourceName, "homepage", fmt.Sprintf("http://homepage-%s", updatedId)),
				),
			},
		},
	})
}

func TestAccAzureRMAdApplication_keyCredential(t *testing.T) {
	resourceName := "azurerm_ad_application.test"
	id := uuid.New().String()
	keyId := uuid.New().String()
	config := testAccAzureRMAdApplication_keyCredential_single(id, keyId, "AsymmetricX509Cert")

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

func TestAccAzureRMAdApplication_updateKeyCredential_changeAttributes(t *testing.T) {
	resourceName := "azurerm_ad_application.test"
	id := uuid.New().String()
	keyId := uuid.New().String()
	keyId2 := uuid.New().String()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProvidersWithTLS,
		CheckDestroy: testCheckAzureRMAdApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAdApplication_keyCredential_single(id, keyId, "AsymmetricX509Cert"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdApplicationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "key_credential.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName, "object_id"),
				),
			},
			{
				Config:      testAccAzureRMAdApplication_keyCredential_single(id, keyId, "Symmetric"),
				ExpectError: regexp.MustCompile(`Error: changing Key Credential properties on existing KeyID`),
			},
			{
				Config: testAccAzureRMAdApplication_keyCredential_single(id, keyId2, "Symmetric"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdApplicationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "key_credential.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName, "object_id"),
				),
			},
		},
	})
}

func TestAccAzureRMAdApplication_updateKeyCredential_changeCount(t *testing.T) {
	resourceName := "azurerm_ad_application.test"
	id := uuid.New().String()
	keyId := uuid.New().String()
	keyId2 := uuid.New().String()
	configSingle := testAccAzureRMAdApplication_keyCredential_single(id, keyId, "AsymmetricX509Cert")
	configDouble := testAccAzureRMAdApplication_keyCredential_double(id, keyId, keyId2)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProvidersWithTLS,
		CheckDestroy: testCheckAzureRMAdApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: configSingle,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdApplicationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", id),
					resource.TestCheckResourceAttr(resourceName, "key_credential.#", "1"),
				),
			},
			{
				Config: configDouble,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdApplicationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", id),
					resource.TestCheckResourceAttr(resourceName, "key_credential.#", "2"),
				),
			},
			{
				Config: configSingle,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdApplicationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", id),
					resource.TestCheckResourceAttr(resourceName, "key_credential.#", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMAdApplication_passwordCredential(t *testing.T) {
	resourceName := "azurerm_ad_application.test"
	id := uuid.New().String()
	keyId := uuid.New().String()
	timeStart := time.Now().UTC()
	timeEnd := timeStart.Add(time.Duration(1) * time.Hour)
	config := testAccAzureRMAdApplication_passwordCredential_single(id, keyId, timeStart, timeEnd)

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

func TestAccAzureRMAdApplication_updatePasswordCredential_changeAttributes(t *testing.T) {
	resourceName := "azurerm_ad_application.test"
	id := uuid.New().String()
	keyId := uuid.New().String()
	keyId2 := uuid.New().String()
	timeStart := time.Now().UTC()
	timeEnd := timeStart.Add(time.Duration(1) * time.Hour)
	timeEnd2 := timeEnd.Add(time.Duration(1) * time.Hour)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAdApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAdApplication_passwordCredential_single(id, keyId, timeStart, timeEnd),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdApplicationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "password_credential.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName, "object_id"),
				),
			},
			{
				Config:      testAccAzureRMAdApplication_passwordCredential_single(id, keyId, timeStart, timeEnd2),
				ExpectError: regexp.MustCompile(`Error: changing Password Credential properties on existing KeyID`),
			},
			{
				Config: testAccAzureRMAdApplication_passwordCredential_single(id, keyId2, timeStart, timeEnd2),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdApplicationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "password_credential.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "app_id"),
					resource.TestCheckResourceAttrSet(resourceName, "object_id"),
				),
			},
		},
	})
}

func TestAccAzureRMAdApplication_updatePasswordCredential_changeCount(t *testing.T) {
	resourceName := "azurerm_ad_application.test"
	id := uuid.New().String()
	keyId := uuid.New().String()
	keyId2 := uuid.New().String()
	timeStart := time.Now().UTC()
	timeEnd := timeStart.Add(time.Duration(1) * time.Hour)
	configSingle := testAccAzureRMAdApplication_passwordCredential_single(id, keyId, timeStart, timeEnd)
	configDouble := testAccAzureRMAdApplication_passwordCredential_double(id, keyId, keyId2, timeStart, timeEnd)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAdApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: configSingle,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdApplicationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "password_credential.#", "1"),
				),
			},
			{
				Config: configDouble,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdApplicationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "password_credential.#", "2"),
				),
			},
			{
				Config: configSingle,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAdApplicationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "password_credential.#", "1"),
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

func testAccAzureRMAdApplication_keyCredential_single(id string, keyId string, usage string) string {
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
    type = "%s"
    usage = "Verify"
    start_date = "${tls_self_signed_cert.test.validity_start_time}"
    end_date = "${tls_self_signed_cert.test.validity_end_time}"
    value = "${tls_self_signed_cert.test.cert_pem}"
  }
}
`, id, keyId, usage)
}

func testAccAzureRMAdApplication_keyCredential_double(id string, keyId string, keyId2 string) string {
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

resource "tls_private_key" "test2" {
  algorithm   = "ECDSA"
  ecdsa_curve = "P384"
}

resource "tls_self_signed_cert" "test2" {
  key_algorithm   = "${tls_private_key.test2.algorithm}"
  private_key_pem = "${tls_private_key.test2.private_key_pem}"

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
    start_date = "${tls_self_signed_cert.test.validity_start_time}"
    end_date = "${tls_self_signed_cert.test.validity_end_time}"
    value = "${tls_self_signed_cert.test.cert_pem}"
  }

  key_credential {
    key_id = "%s"
    type = "AsymmetricX509Cert"
    usage = "Verify"
    start_date = "${tls_self_signed_cert.test2.validity_start_time}"
    end_date = "${tls_self_signed_cert.test2.validity_end_time}"
    value = "${tls_self_signed_cert.test2.cert_pem}"
  }
}
`, id, keyId, keyId2)
}

func testAccAzureRMAdApplication_passwordCredential_single(id string, keyId string, timeStart time.Time, timeEnd time.Time) string {
	ts := string(timeStart.Format(time.RFC3339))
	te := string(timeEnd.Format(time.RFC3339))
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
`, id, keyId, ts, te)
}

func testAccAzureRMAdApplication_passwordCredential_double(id string, keyId string, keyId2 string, timeStart time.Time, timeEnd time.Time) string {
	ts := string(timeStart.Format(time.RFC3339))
	te := string(timeEnd.Format(time.RFC3339))
	return fmt.Sprintf(`
resource "azurerm_ad_application" "test" {
  display_name = "%s"

  password_credential {
    key_id = "%s"
    value = "test"
    start_date = "%s"
    end_date = "%s"
  }

  password_credential {
    key_id = "%s"
    value = "test"
    start_date = "%s"
    end_date = "%s"
  }
}
`, id, keyId, ts, te, keyId2, ts, te)
}
