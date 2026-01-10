// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CdnFrontDoorCustomDomainResource struct{}

func TestAccCdnFrontDoorCustomDomain_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_custom_domain", "test")
	r := CdnFrontDoorCustomDomainResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorCustomDomain_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_custom_domain", "test")
	r := CdnFrontDoorCustomDomainResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccCdnFrontDoorCustomDomain_update(t *testing.T) {
	if features.FivePointOh() {
		t.Skipf("There is no available `tls_version` to test update, to test CMK, it requires an official certificate from approved provider list instead of testing cert.")
	}
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_custom_domain", "test")
	r := CdnFrontDoorCustomDomainResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorCustomDomain_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_custom_domain", "test")
	r := CdnFrontDoorCustomDomainResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorCustomDomain_tlsVersion_legacy(t *testing.T) {
	if features.FivePointOh() {
		t.Skip("Skipping test in 5.0 mode as `minimum_tls_version` field was removed")
	}
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_custom_domain", "test")
	r := CdnFrontDoorCustomDomainResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.tlsVersionLegacy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontDoorCustomDomain_cipherSuites_validation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_custom_domain", "test")
	r := CdnFrontDoorCustomDomainResource{}

	testSteps := []acceptance.TestStep{
		{
			Config:      r.customizedCipherSuiteWithoutBlock(data),
			ExpectError: regexp.MustCompile("`custom_ciphers` is required when `type` is `Customized`"),
		},
		{
			Config:      r.customizedCipherSuiteEmpty(data),
			ExpectError: regexp.MustCompile("at least one cipher suite must be selected in `custom_ciphers` when `type` is set to `Customized`"),
		},
		{
			Config:      r.customizedCipherSuiteTls12MissingWithTls12Min(data),
			ExpectError: regexp.MustCompile("at least one TLS 1.2 cipher suite must be specified in `custom_ciphers.tls12` when `minimum_version` is set to `TLS12`"),
		},
		{
			Config:      r.customCiphersWithPresetType(data),
			ExpectError: regexp.MustCompile("`custom_ciphers` cannot be specified when `type` is not `Customized`"),
		},
	}

	if !features.FivePointOh() {
		testSteps = append(testSteps, acceptance.TestStep{
			Config:      r.deprecatedMinimumTlsVersionWithTls10(data),
			ExpectError: regexp.MustCompile("support for TLS 1.0 and 1.1 was retired on March 1, 2025"),
		})
	}

	data.ResourceTest(t, r, testSteps)
}

func TestAccCdnFrontDoorCustomDomain_cipherSuites_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_custom_domain", "test")
	r := CdnFrontDoorCustomDomainResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.cipherSuitesTls12Single(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.cipherSuitesTls12Multiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.cipherSuitesMixedWithTls12MinSingle(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.cipherSuitesMixedWithTls12MinMultiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r CdnFrontDoorCustomDomainResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FrontDoorCustomDomainID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cdn.FrontDoorCustomDomainsClient
	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.CustomDomainName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(true), nil
}

func (r CdnFrontDoorCustomDomainResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_custom_domain" "test" {
  name                     = "acctestcustomdomain-%d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
  dns_zone_id              = azurerm_dns_zone.test.id
  host_name                = join(".", ["%s", azurerm_dns_zone.test.name])

  tls {
    certificate_type = "ManagedCertificate"
    minimum_version  = "TLS12"
  }
}
`, template, data.RandomInteger, data.RandomString)
}

func (r CdnFrontDoorCustomDomainResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_custom_domain" "import" {
  name                     = azurerm_cdn_frontdoor_custom_domain.test.name
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_custom_domain.test.cdn_frontdoor_profile_id
  dns_zone_id              = azurerm_cdn_frontdoor_custom_domain.test.dns_zone_id
  host_name                = azurerm_cdn_frontdoor_custom_domain.test.host_name

  tls {
    certificate_type = "ManagedCertificate"
    minimum_version  = "TLS12"
  }
}
`, config)
}

func (r CdnFrontDoorCustomDomainResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_custom_domain" "test" {
  name                     = "acctestcustomdomain-%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
  dns_zone_id              = azurerm_dns_zone.test.id
  host_name                = join(".", ["%s", azurerm_dns_zone.test.name])

  tls {
    certificate_type    = "ManagedCertificate"
    minimum_tls_version = "TLS10"
  }
}
`, template, data.RandomInteger, data.RandomString)
}

func (r CdnFrontDoorCustomDomainResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_custom_domain" "test" {
  name                     = "acctestcustomdomain-%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
  dns_zone_id              = azurerm_dns_zone.test.id
  host_name                = join(".", ["%s", azurerm_dns_zone.test.name])

  tls {
    certificate_type = "ManagedCertificate"
    minimum_version  = "TLS12"
  }

}
`, template, data.RandomInteger, data.RandomString)
}

func (r CdnFrontDoorCustomDomainResource) tlsVersionLegacy(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_custom_domain" "test" {
  name                     = "acctestcustomdomain-%d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
  dns_zone_id              = azurerm_dns_zone.test.id
  host_name                = join(".", ["%s", azurerm_dns_zone.test.name])

  tls {
    certificate_type    = "ManagedCertificate"
    minimum_tls_version = "TLS12"
  }
}
`, template, data.RandomInteger, data.RandomString)
}

// TODO: Add test case that uses pre_validated_custom_domain_resource_id
// TODO: Add test case that uses CMK, this cannot be a test cert or a self
// signed cert it must be an official cert from the approved list of cert
// providers by the service.

func (r CdnFrontDoorCustomDomainResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-afdx-%[1]d"
  location = "%[2]s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%[1]d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctestcdnfdprofile-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Standard_AzureFrontDoor"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r CdnFrontDoorCustomDomainResource) customizedCipherSuiteWithoutBlock(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cdn_frontdoor_custom_domain" "test" {
  name                     = "acctest-customdomain-%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
  dns_zone_id              = azurerm_dns_zone.test.id
  host_name                = "acctest-%[2]d.acctestzone%[2]d.com"

  tls {
    certificate_type = "ManagedCertificate"
    minimum_version  = "TLS12"

    cipher_suite {
      type = "Customized"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorCustomDomainResource) customizedCipherSuiteEmpty(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cdn_frontdoor_custom_domain" "test" {
  name                     = "acctest-customdomain-%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
  dns_zone_id              = azurerm_dns_zone.test.id
  host_name                = "acctest-%[2]d.acctestzone%[2]d.com"

  tls {
    certificate_type = "ManagedCertificate"
    minimum_version  = "TLS12"

    cipher_suite {
      type = "Customized"

      custom_ciphers {
        # Empty - no cipher suites selected
      }
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorCustomDomainResource) customizedCipherSuiteTls12MissingWithTls12Min(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cdn_frontdoor_custom_domain" "test" {
  name                     = "acctest-customdomain-%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
  dns_zone_id              = azurerm_dns_zone.test.id
  host_name                = "acctest-%[2]d.acctestzone%[2]d.com"

  tls {
    certificate_type = "ManagedCertificate"
    minimum_version  = "TLS12"

    cipher_suite {
      type = "Customized"

      custom_ciphers {
        # Invalid: minimum_version is TLS12 but only TLS13 ciphers defined
        tls13 = [
          "TLS_AES_256_GCM_SHA384",
        ]
      }
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorCustomDomainResource) cipherSuitesTls12Single(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cdn_frontdoor_custom_domain" "test" {
  name                     = "acctestcustomdomain-%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
  dns_zone_id              = azurerm_dns_zone.test.id
  host_name                = join(".", ["%[3]s", azurerm_dns_zone.test.name])

  tls {
    certificate_type = "ManagedCertificate"
    minimum_version  = "TLS12"

    cipher_suite {
      type = "Customized"

      custom_ciphers {
        tls12 = [
          "ECDHE_RSA_AES256_GCM_SHA384",
        ]
      }
    }
  }
}
`, r.template(data), data.RandomInteger, data.RandomString)
}

func (r CdnFrontDoorCustomDomainResource) cipherSuitesTls12Multiple(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cdn_frontdoor_custom_domain" "test" {
  name                     = "acctestcustomdomain-%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
  dns_zone_id              = azurerm_dns_zone.test.id
  host_name                = join(".", ["%[3]s", azurerm_dns_zone.test.name])

  tls {
    certificate_type = "ManagedCertificate"
    minimum_version  = "TLS12"

    cipher_suite {
      type = "Customized"

      custom_ciphers {
        tls12 = [
          "ECDHE_RSA_AES128_GCM_SHA256",
          "ECDHE_RSA_AES256_GCM_SHA384",
          "DHE_RSA_AES128_GCM_SHA256",
        ]
      }
    }
  }
}
`, r.template(data), data.RandomInteger, data.RandomString)
}

func (r CdnFrontDoorCustomDomainResource) cipherSuitesMixedWithTls12MinSingle(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cdn_frontdoor_custom_domain" "test" {
  name                     = "acctestcustomdomain-%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
  dns_zone_id              = azurerm_dns_zone.test.id
  host_name                = join(".", ["%[3]s", azurerm_dns_zone.test.name])

  tls {
    certificate_type = "ManagedCertificate"
    minimum_version  = "TLS12"

    cipher_suite {
      type = "Customized"

      custom_ciphers {
        tls12 = [
          "ECDHE_RSA_AES128_GCM_SHA256",
        ]
        tls13 = [
          "TLS_AES_256_GCM_SHA384",
        ]
      }
    }
  }
}
`, r.template(data), data.RandomInteger, data.RandomString)
}

func (r CdnFrontDoorCustomDomainResource) cipherSuitesMixedWithTls12MinMultiple(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cdn_frontdoor_custom_domain" "test" {
  name                     = "acctestcustomdomain-%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
  dns_zone_id              = azurerm_dns_zone.test.id
  host_name                = join(".", ["%[3]s", azurerm_dns_zone.test.name])

  tls {
    certificate_type = "ManagedCertificate"
    minimum_version  = "TLS12"

    cipher_suite {
      type = "Customized"

      custom_ciphers {
        tls12 = [
          "ECDHE_RSA_AES128_GCM_SHA256",
        ]
        tls13 = [
          "TLS_AES_128_GCM_SHA256",
          "TLS_AES_256_GCM_SHA384",
        ]
      }
    }
  }
}
`, r.template(data), data.RandomInteger, data.RandomString)
}

func (r CdnFrontDoorCustomDomainResource) customCiphersWithPresetType(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cdn_frontdoor_custom_domain" "test" {
  name                     = "acctest-customdomain-%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
  dns_zone_id              = azurerm_dns_zone.test.id
  host_name                = "acctest-%[2]d.acctestzone%[2]d.com"

  tls {
    certificate_type = "ManagedCertificate"
    minimum_version  = "TLS12"

    cipher_suite {
      type = "TLS12_2023"

      custom_ciphers {
        # Invalid: custom_ciphers cannot be specified when type is not Customized
        tls12 = [
          "ECDHE_RSA_AES128_GCM_SHA256",
        ]
      }
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontDoorCustomDomainResource) deprecatedMinimumTlsVersionWithTls10(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cdn_frontdoor_custom_domain" "test" {
  name                     = "acctest-customdomain-%[2]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
  dns_zone_id              = azurerm_dns_zone.test.id
  host_name                = "acctest-%[2]d.acctestzone%[2]d.com"

  tls {
    certificate_type    = "ManagedCertificate"
    minimum_tls_version = "TLS10"
  }
}
`, r.template(data), data.RandomInteger)
}
