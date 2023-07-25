// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CdnEndpointCustomDomainResource struct {
	DNSZoneRG   string
	DNSZoneName string

	// Subdomain Name, this needs to be the same value as set to "CN" in the certificate.
	SubDomainName string

	// PFX File (with base64 encoded)
	CertificateP12 string
}

func NewCdnEndpointCustomDomainResource(dnsZoneRg, dnsZoneName string) *CdnEndpointCustomDomainResource {
	return &CdnEndpointCustomDomainResource{
		DNSZoneRG:     dnsZoneRg,
		DNSZoneName:   dnsZoneName,
		SubDomainName: acceptance.RandString(3),
	}
}

func TestAccCdnEndpointCustomDomain_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint_custom_domain", "test")

	r := NewCdnEndpointCustomDomainResource(os.Getenv("ARM_TEST_DNS_ZONE_RESOURCE_GROUP_NAME"), os.Getenv("ARM_TEST_DNS_ZONE_NAME"))
	r.preCheck(t)

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

func TestAccCdnEndpointCustomDomain_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint_custom_domain", "test")

	r := NewCdnEndpointCustomDomainResource(os.Getenv("ARM_TEST_DNS_ZONE_RESOURCE_GROUP_NAME"), os.Getenv("ARM_TEST_DNS_ZONE_NAME"))
	r.preCheck(t)

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

func TestAccCdnEndpointCustomDomain_httpsCdn(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint_custom_domain", "test")

	r := NewCdnEndpointCustomDomainResource(os.Getenv("ARM_TEST_DNS_ZONE_RESOURCE_GROUP_NAME"), os.Getenv("ARM_TEST_DNS_ZONE_NAME"))
	r.preCheck(t)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.httpsCdn(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnEndpointCustomDomain_httpsUserManagedCertificate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint_custom_domain", "test")

	r := NewCdnEndpointCustomDomainResource(os.Getenv("ARM_TEST_DNS_ZONE_RESOURCE_GROUP_NAME"), os.Getenv("ARM_TEST_DNS_ZONE_NAME"))
	r.CertificateP12 = os.Getenv("ARM_TEST_DNS_CERTIFICATE")
	r.SubDomainName = os.Getenv("ARM_TEST_DNS_SUBDOMAIN_NAME")
	r.preCheckUserManagedCertificate(t)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.httpsUserManagedCertificate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// The "key_vault_secret_id" is skipped here since during import, there is no knowledge about whether users want
		// versioned or versionless certificate id. That means the imported "key_vault_secret_id" is what it is at the
		// remote API representation, which might be different than it as defined in the configuration.
		data.ImportStep("user_managed_https.0.key_vault_secret_id", "user_managed_https.0.key_vault_certificate_id"),
	})
}

func TestAccCdnEndpointCustomDomain_httpsUserManagedCertificateDeprecated(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skipf("This test is skipped since v4.0")
	}
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint_custom_domain", "test")

	r := NewCdnEndpointCustomDomainResource(os.Getenv("ARM_TEST_DNS_ZONE_RESOURCE_GROUP_NAME"), os.Getenv("ARM_TEST_DNS_ZONE_NAME"))
	r.CertificateP12 = os.Getenv("ARM_TEST_DNS_CERTIFICATE")
	r.SubDomainName = os.Getenv("ARM_TEST_DNS_SUBDOMAIN_NAME")
	r.preCheckUserManagedCertificate(t)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.httpsUserManagedCertificateDeprecated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("user_managed_https.0.key_vault_secret_id", "user_managed_https.0.key_vault_certificate_id"),
	})
}

func TestAccCdnEndpointCustomDomain_httpsUserManagedSecret(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint_custom_domain", "test")

	r := NewCdnEndpointCustomDomainResource(os.Getenv("ARM_TEST_DNS_ZONE_RESOURCE_GROUP_NAME"), os.Getenv("ARM_TEST_DNS_ZONE_NAME"))
	r.CertificateP12 = os.Getenv("ARM_TEST_DNS_CERTIFICATE")
	r.SubDomainName = os.Getenv("ARM_TEST_DNS_SUBDOMAIN_NAME")
	r.preCheckUserManagedCertificate(t)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.httpsUserManagedSecret(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("user_managed_https.0.key_vault_secret_id", "user_managed_https.0.key_vault_certificate_id"),
	})
}

func TestAccCdnEndpointCustomDomain_httpsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint_custom_domain", "test")

	r := NewCdnEndpointCustomDomainResource(os.Getenv("ARM_TEST_DNS_ZONE_RESOURCE_GROUP_NAME"), os.Getenv("ARM_TEST_DNS_ZONE_NAME"))
	r.CertificateP12 = os.Getenv("ARM_TEST_DNS_CERTIFICATE")
	r.SubDomainName = os.Getenv("ARM_TEST_DNS_SUBDOMAIN_NAME")
	r.preCheckUserManagedCertificate(t)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.httpsCdn(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.httpsUserManagedCertificate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("user_managed_https.0.key_vault_secret_id", "user_managed_https.0.key_vault_certificate_id"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r CdnEndpointCustomDomainResource) preCheck(t *testing.T) {
	if r.DNSZoneRG == "" {
		t.Skipf("`ARM_TEST_DNS_ZONE_RESOURCE_GROUP_NAME` must be set for acceptance tests!")
	}
	if r.DNSZoneName == "" {
		t.Skipf("`ARM_TEST_DNS_ZONE_NAME` must be set for acceptance tests!")
	}
}

func (r CdnEndpointCustomDomainResource) preCheckUserManagedCertificate(t *testing.T) {
	r.preCheck(t)
	if r.SubDomainName == "" {
		t.Skipf("`ARM_TEST_DNS_SUBDOMAIN_NAME` must be set for acceptance tests!")
	}
	if r.CertificateP12 == "" {
		t.Skipf("`ARM_TEST_DNS_CERTIFICATE` must be set for acceptance tests!")
	}
}

func (r CdnEndpointCustomDomainResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.CustomDomainID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Cdn.CustomDomainsClient.Get(ctx, id.ResourceGroup, id.ProfileName, id.EndpointName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %q: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r CdnEndpointCustomDomainResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.CustomDomainID(state.ID)
	if err != nil {
		return nil, err
	}

	c := client.Cdn.CustomDomainsClient
	future, err := c.Delete(ctx, id.ResourceGroup, id.ProfileName, id.EndpointName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("deleting %q: %+v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, c.Client); err != nil {
		return nil, fmt.Errorf("waiting for deletion of %q: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (r CdnEndpointCustomDomainResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_endpoint_custom_domain" "test" {
  name            = "acceptance-customdomain"
  cdn_endpoint_id = azurerm_cdn_endpoint.test.id
  host_name       = "${azurerm_dns_cname_record.test.name}.${data.azurerm_dns_zone.test.name}"
}
`, template)
}

func (r CdnEndpointCustomDomainResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_endpoint_custom_domain" "import" {
  name            = azurerm_cdn_endpoint_custom_domain.test.name
  cdn_endpoint_id = azurerm_cdn_endpoint_custom_domain.test.cdn_endpoint_id
  host_name       = azurerm_cdn_endpoint_custom_domain.test.host_name
}
`, template)
}

func (r CdnEndpointCustomDomainResource) httpsCdn(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_cdn_endpoint_custom_domain" "test" {
  name            = "testcustomdomain-%[2]d"
  cdn_endpoint_id = azurerm_cdn_endpoint.test.id
  host_name       = "${azurerm_dns_cname_record.test.name}.${data.azurerm_dns_zone.test.name}"
  cdn_managed_https {
    certificate_type = "Dedicated"
    protocol_type    = "ServerNameIndication"
  }
}
`, template, data.RandomIntOfLength(8))
}

func (r CdnEndpointCustomDomainResource) httpsUserManagedBase(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

data "azurerm_client_config" "test" {
}

data "azuread_service_principal" "test" {
  display_name = "Microsoft.Azure.Cdn"
}

resource "azurerm_key_vault" "test" {
  name                = "testkeyvault-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.test.tenant_id
  sku_name            = "standard"
  access_policy {
    tenant_id = data.azurerm_client_config.test.tenant_id
    object_id = data.azurerm_client_config.test.object_id
    certificate_permissions = [
      "Get",
      "Delete",
      "Import",
      "Purge",
    ]
    key_permissions = [
      "Get",
      "Create",
    ]
    secret_permissions = [
      "Get",
      "Set",
      "Delete",
      "Purge",
    ]
  }
  access_policy {
    tenant_id = data.azurerm_client_config.test.tenant_id
    object_id = data.azuread_service_principal.test.object_id
    certificate_permissions = [
      "List",
      "Get",
    ]
    secret_permissions = [
      "List",
      "Get",
    ]
  }
}
`, template, data.RandomIntOfLength(8))
}

func (r CdnEndpointCustomDomainResource) httpsUserManagedCertificate(data acceptance.TestData) string {
	template := r.httpsUserManagedBase(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_key_vault_certificate" "test" {
  name         = "testkeyvaultcert-%[2]d"
  key_vault_id = azurerm_key_vault.test.id
  certificate {
    contents = file("%[3]s")
    password = ""
  }
  certificate_policy {
    issuer_parameters {
      name = "Self"
    }
    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = false
    }
    secret_properties {
      content_type = "application/x-pkcs12"
    }
  }
}
resource "azurerm_cdn_endpoint_custom_domain" "test" {
  name            = "testcustomdomain-%[2]d"
  cdn_endpoint_id = azurerm_cdn_endpoint.test.id
  host_name       = "${azurerm_dns_cname_record.test.name}.${data.azurerm_dns_zone.test.name}"
  user_managed_https {
    key_vault_secret_id = azurerm_key_vault_certificate.test.secret_id
  }
}
`, template, data.RandomIntOfLength(8), r.CertificateP12)
}

func (r CdnEndpointCustomDomainResource) httpsUserManagedCertificateDeprecated(data acceptance.TestData) string {
	template := r.httpsUserManagedBase(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_key_vault_certificate" "test" {
  name         = "testkeyvaultcert-%[2]d"
  key_vault_id = azurerm_key_vault.test.id
  certificate {
    contents = file("%[3]s")
    password = ""
  }
  certificate_policy {
    issuer_parameters {
      name = "Self"
    }
    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = false
    }
    secret_properties {
      content_type = "application/x-pkcs12"
    }
  }
}
resource "azurerm_cdn_endpoint_custom_domain" "test" {
  name            = "testcustomdomain-%[2]d"
  cdn_endpoint_id = azurerm_cdn_endpoint.test.id
  host_name       = "${azurerm_dns_cname_record.test.name}.${data.azurerm_dns_zone.test.name}"
  user_managed_https {
    key_vault_certificate_id = azurerm_key_vault_certificate.test.id
  }
}
`, template, data.RandomIntOfLength(8), r.CertificateP12)
}

func (r CdnEndpointCustomDomainResource) httpsUserManagedSecret(data acceptance.TestData) string {
	template := r.httpsUserManagedBase(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_key_vault_secret" "test" {
  name         = "testkeyvaultsecret-%[2]d"
  key_vault_id = azurerm_key_vault.test.id
  content_type = "application/x-pkcs12"
  value        = file("%[3]s")
}
resource "azurerm_cdn_endpoint_custom_domain" "test" {
  name            = "testcustomdomain-%[2]d"
  cdn_endpoint_id = azurerm_cdn_endpoint.test.id
  host_name       = "${azurerm_dns_cname_record.test.name}.${data.azurerm_dns_zone.test.name}"
  user_managed_https {
    key_vault_secret_id = azurerm_key_vault_secret.test.id
  }
}
`, template, data.RandomIntOfLength(8), r.CertificateP12)
}

func (r CdnEndpointCustomDomainResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acceptanceRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acceptancesa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acceptancecdnprof%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Microsoft"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acceptancecdnend%[1]d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  origin {
    name      = "test"
    host_name = azurerm_storage_account.test.primary_blob_host
  }
}

data "azurerm_dns_zone" "test" {
  name                = "%[3]s"
  resource_group_name = "%[4]s"
}

resource "azurerm_dns_cname_record" "test" {
  name                = "%[5]s"
  zone_name           = data.azurerm_dns_zone.test.name
  resource_group_name = data.azurerm_dns_zone.test.resource_group_name
  ttl                 = 3600
  target_resource_id  = azurerm_cdn_endpoint.test.id
}
`, data.RandomIntOfLength(8), data.Locations.Primary, r.DNSZoneName, r.DNSZoneRG, r.SubDomainName)
}
