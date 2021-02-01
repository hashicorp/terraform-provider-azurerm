package springcloud_test

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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/springcloud/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type SpringCloudCustomDomainResource struct {
}

func TestAccSpringCloudCustomDomain_basic(t *testing.T) {
	if os.Getenv("ARM_TEST_DNS_ZONE") == "" || os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" {
		t.Skip("Skipping as ARM_TEST_DNS_ZONE and/or ARM_TEST_DATA_RESOURCE_GROUP are not specified")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_custom_domain", "test")
	r := SpringCloudCustomDomainResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSpringCloudCustomDomain_requiresImport(t *testing.T) {
	if os.Getenv("ARM_TEST_DNS_ZONE") == "" || os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" {
		t.Skip("Skipping as ARM_TEST_DNS_ZONE and/or ARM_TEST_DATA_RESOURCE_GROUP are not specified")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_custom_domain", "test")
	r := SpringCloudCustomDomainResource{}

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

func TestAccSpringCloudCustomDomain_complete(t *testing.T) {
	if os.Getenv("ARM_TEST_DNS_ZONE") == "" || os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" {
		t.Skip("Skipping as ARM_TEST_DNS_ZONE and/or ARM_TEST_DATA_RESOURCE_GROUP are not specified")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_custom_domain", "test")
	r := SpringCloudCustomDomainResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSpringCloudCustomDomain_update(t *testing.T) {
	if os.Getenv("ARM_TEST_DNS_ZONE") == "" || os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" {
		t.Skip("Skipping as ARM_TEST_DNS_ZONE and/or ARM_TEST_DATA_RESOURCE_GROUP are not specified")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_custom_domain", "test")
	r := SpringCloudJavaDeploymentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r SpringCloudCustomDomainResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.SpringCloudCustomDomainID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.AppPlatform.CustomDomainsClient.Get(ctx, id.ResourceGroup, id.SpringName, id.AppName, id.DomainName)
	if err != nil {
		return nil, fmt.Errorf("reading Spring Cloud Custom Domain %q (Spring Cloud Service %q / App %q / Resource Group %q): %+v", id.DomainName, id.SpringName, id.AppName, id.ResourceGroup, err)
	}

	return utils.Bool(resp.Properties != nil), nil
}

func (r SpringCloudCustomDomainResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_custom_domain" "test" {
  name                = join(".", [azurerm_dns_cname_record.test.name, azurerm_dns_cname_record.test.zone_name])
  spring_cloud_app_id = azurerm_spring_cloud_app.test.id
}
`, r.template(data))
}

func (r SpringCloudCustomDomainResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_custom_domain" "import" {
  name                = azurerm_spring_cloud_custom_domain.test.name
  spring_cloud_app_id = azurerm_spring_cloud_custom_domain.test.spring_cloud_app_id
}
`, r.basic(data))
}

func (r SpringCloudCustomDomainResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {
}

data "azuread_service_principal" "test" {
  display_name = "Azure Spring Cloud Domain-Management"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkeyvault%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"

  access_policy {
    tenant_id               = data.azurerm_client_config.current.tenant_id
    object_id               = data.azurerm_client_config.current.object_id
    secret_permissions      = ["set"]
    certificate_permissions = ["create", "delete", "get", "update"]
  }

  access_policy {
    tenant_id               = data.azurerm_client_config.current.tenant_id
    object_id               = data.azuread_service_principal.test.object_id
    secret_permissions      = ["get", "list"]
    certificate_permissions = ["get", "list"]
  }
}

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctestcert%s"
  key_vault_id = azurerm_key_vault.test.id

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }

    lifetime_action {
      action {
        action_type = "AutoRenew"
      }

      trigger {
        days_before_expiry = 30
      }
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    x509_certificate_properties {
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]

      subject            = join("", ["CN=*.", data.azurerm_dns_zone.test.name])

      subject_alternative_names {
        dns_names = [join(".", [azurerm_dns_cname_record.test.name, azurerm_dns_cname_record.test.zone_name])]
      }

      validity_in_months = 12
    }
  }
}

resource "azurerm_spring_cloud_certificate" "test" {
  name                     = "acctest-scc-%d"
  resource_group_name      = azurerm_spring_cloud_service.test.resource_group_name
  service_name             = azurerm_spring_cloud_service.test.name
  key_vault_certificate_id = azurerm_key_vault_certificate.test.id
}

resource "azurerm_spring_cloud_custom_domain" "test" {
  name                = join(".", [azurerm_dns_cname_record.test.name, azurerm_dns_cname_record.test.zone_name])
  spring_cloud_app_id = azurerm_spring_cloud_custom_domain.test.spring_cloud_app_id
  cert_name           = azurerm_spring_cloud_certificate.test.name
  thumbprint          = azurerm_spring_cloud_certificate.test.thumbprint
}
`, r.template(data), data.RandomString, data.RandomString, data.RandomInteger)
}

func (SpringCloudCustomDomainResource) template(data acceptance.TestData) string {
	dnsZone := os.Getenv("ARM_TEST_DNS_ZONE")
	dataResourceGroup := os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-spring-%d"
  location = "%s"
}

resource "azurerm_spring_cloud_service" "test" {
  name                = "acctest-sc-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_spring_cloud_app" "test" {
  name                = "acctest-sca-%d"
  resource_group_name = azurerm_spring_cloud_service.test.resource_group_name
  service_name        = azurerm_spring_cloud_service.test.name
}

data "azurerm_dns_zone" "test" {
  name                = "%s"
  resource_group_name = "%s"
}

resource "azurerm_dns_cname_record" "test" {
  name                = "%s"
  zone_name           = data.azurerm_dns_zone.test.name
  resource_group_name = data.azurerm_dns_zone.test.resource_group_name
  ttl                 = 300
  record              = azurerm_spring_cloud_app.test.fqdn
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, dnsZone, dataResourceGroup, data.RandomString)
}
