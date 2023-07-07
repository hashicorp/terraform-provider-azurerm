// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logic_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountcertificates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LogicAppIntegrationAccountCertificateResource struct{}

func TestAccLogicAppIntegrationAccountCertificate_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account_certificate", "test")
	r := LogicAppIntegrationAccountCertificateResource{}
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

func TestAccLogicAppIntegrationAccountCertificate_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account_certificate", "test")
	r := LogicAppIntegrationAccountCertificateResource{}

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

func TestAccLogicAppIntegrationAccountCertificate_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account_certificate", "test")
	r := LogicAppIntegrationAccountCertificateResource{}

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

func TestAccLogicAppIntegrationAccountCertificate_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account_certificate", "test")
	r := LogicAppIntegrationAccountCertificateResource{}

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

func (r LogicAppIntegrationAccountCertificateResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := integrationaccountcertificates.ParseCertificateID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Logic.IntegrationAccountCertificateClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %q %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (r LogicAppIntegrationAccountCertificateResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

data "azuread_service_principal" "test" {
  display_name = "Azure Logic Apps"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-logic-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"
}

resource "azurerm_key_vault_access_policy" "test1" {
  key_vault_id = azurerm_key_vault.test.id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "Purge",
    "Recover",
    "Update",
    "List",
    "Decrypt",
    "Sign",
    "GetRotationPolicy"
  ]

  secret_permissions = [
    "Get",
    "Set",
  ]

  tenant_id = data.azurerm_client_config.current.tenant_id
  object_id = data.azuread_service_principal.test.object_id
}

resource "azurerm_key_vault_access_policy" "test2" {
  key_vault_id = azurerm_key_vault.test.id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "Purge",
    "Recover",
    "Update",
    "List",
    "Decrypt",
    "Sign",
    "GetRotationPolicy"
  ]

  secret_permissions = [
    "Get",
    "Set",
  ]

  tenant_id = data.azurerm_client_config.current.tenant_id
  object_id = data.azurerm_client_config.current.object_id
}

resource "azurerm_key_vault_key" "test" {
  name         = "acctestkvkey-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey"
  ]

  depends_on = [azurerm_key_vault_access_policy.test1, azurerm_key_vault_access_policy.test2]
}

resource "azurerm_logic_app_integration_account" "test" {
  name                = "acctestia-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Standard"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomInteger)
}

func (r LogicAppIntegrationAccountCertificateResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_integration_account_certificate" "test" {
  name                     = "acctest-iac-%d"
  resource_group_name      = azurerm_resource_group.test.name
  integration_account_name = azurerm_logic_app_integration_account.test.name

  key_vault_key {
    key_name     = azurerm_key_vault_key.test.name
    key_vault_id = azurerm_key_vault.test.id
    key_version  = azurerm_key_vault_key.test.version
  }
}
`, template, data.RandomInteger)
}

func (r LogicAppIntegrationAccountCertificateResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_integration_account_certificate" "import" {
  name                     = azurerm_logic_app_integration_account_certificate.test.name
  resource_group_name      = azurerm_logic_app_integration_account_certificate.test.resource_group_name
  integration_account_name = azurerm_logic_app_integration_account_certificate.test.integration_account_name

  key_vault_key {
    key_name     = azurerm_logic_app_integration_account_certificate.test.key_vault_key.0.key_name
    key_vault_id = azurerm_logic_app_integration_account_certificate.test.key_vault_key.0.key_vault_id
    key_version  = azurerm_logic_app_integration_account_certificate.test.key_vault_key.0.key_version
  }
}
`, config)
}

func (r LogicAppIntegrationAccountCertificateResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_integration_account_certificate" "test" {
  name                     = "acctest-iac-%d"
  resource_group_name      = azurerm_resource_group.test.name
  integration_account_name = azurerm_logic_app_integration_account.test.name
  public_certificate       = "MIICsjCCAZoCCQCMdt7DvygPtDANBgkqhkiG9w0BAQsFADAbMRkwFwYDVQQDDBBhcGkudGVycmFmb3JtLmlvMB4XDTE4MDcwNTEwMzMzMFoXDTI4MDcwMjEwMzMzMFowGzEZMBcGA1UEAwwQYXBpLnRlcnJhZm9ybS5pbzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAKQW332Ol28CsidAheD1aL9Ul8JWnKLdaVxKZ3ssl5CXjPDOmM7IXk0SgbQnUC8lIlPFZiDGbQ1sB6OTMun6ZZ4ipLp80dtl0roCLtCnDQOBGzCNArCYAoXRurjkXEY7tpD0wwtU72+37h3HQ4g0VS6VItJCqJ9QADV+HO2ZWuZTez70MhoL6OLfZP7HGYdJDKgfEVNF5XlbVzNAGkDIJFdhjNxyGGu5Nfsm1pfQhAyunkk7JVamjUg5IojRdo63IS9wwzMOdeGSAbBcsJfYeCfVg2kupR8q0TmZ+x93RmmOlbSi66kEYxRzZ9YCQeHJmn1YfJ92BpCUiy9A6Z1iaKUCAwEAATANBgkqhkiG9w0BAQsFAAOCAQEAJ7JhlecP7J48wI2QHTMbAMkkWBv/iWq1/QIF4ugH3Zb5PorOv+NfhQ0LlWiw/SzN8Ae95vUixAGYHMSa28oumM5K1OsqKEkVIo1AoBH8nBz+VcTpRD/mHXotAHPAZt9j5LqeHX+enR6RbINAf3jn+YU3MdVe0MsADdFASVDfjmQP2R7o9aJb/QqOg3bZBWsiBDEISfyaH2+pgUM7wtwEoFWmEMlgjLK1MRBs1cDZXqnHaCd/rs+NmWV9naEu7x5fyQOk4HozkpweR+Jx1sBlTRsa49/qSHt/6ULKfO01/cTs4iF71ykXPbh3Kj9cI2uo9aYtXkxkhKrGyUpA7FJqWw=="

  metadata = <<METADATA
    {
        "foo": "bar"
    }
METADATA

  key_vault_key {
    key_name     = azurerm_key_vault_key.test.name
    key_vault_id = azurerm_key_vault.test.id
    key_version  = azurerm_key_vault_key.test.version
  }
}
`, template, data.RandomInteger)
}

func (r LogicAppIntegrationAccountCertificateResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault" "test2" {
  name                = "acctestkv2-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"
}

resource "azurerm_key_vault_access_policy" "test3" {
  key_vault_id = azurerm_key_vault.test2.id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "Purge",
    "Recover",
    "Update",
    "List",
    "Decrypt",
    "Sign",
    "GetRotationPolicy"
  ]

  secret_permissions = [
    "Get",
    "Set",
  ]

  tenant_id = data.azurerm_client_config.current.tenant_id
  object_id = data.azuread_service_principal.test.object_id
}

resource "azurerm_key_vault_access_policy" "test4" {
  key_vault_id = azurerm_key_vault.test2.id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "Purge",
    "Recover",
    "Update",
    "List",
    "Decrypt",
    "Sign",
    "GetRotationPolicy"
  ]

  secret_permissions = [
    "Get",
    "Set",
  ]

  tenant_id = data.azurerm_client_config.current.tenant_id
  object_id = data.azurerm_client_config.current.object_id
}

resource "azurerm_key_vault_key" "test2" {
  name         = "acctestkvkey2-%s"
  key_vault_id = azurerm_key_vault.test2.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey"
  ]

  depends_on = [azurerm_key_vault_access_policy.test3, azurerm_key_vault_access_policy.test4]
}

resource "azurerm_logic_app_integration_account_certificate" "test" {
  name                     = "acctest-iac-%d"
  resource_group_name      = azurerm_resource_group.test.name
  integration_account_name = azurerm_logic_app_integration_account.test.name
  public_certificate       = "MIIDbzCCAlegAwIBAgIJAIzjRD36sIbbMA0GCSqGSIb3DQEBCwUAME0xCzAJBgNVBAYTAlVTMRMwEQYDVQQIDApTb21lLVN0YXRlMRIwEAYDVQQKDAl0ZXJyYWZvcm0xFTATBgNVBAMMDHRlcnJhZm9ybS5pbzAgFw0xNzA0MjEyMDA1MjdaGA8yMTE3MDMyODIwMDUyN1owTTELMAkGA1UEBhMCVVMxEzARBgNVBAgMClNvbWUtU3RhdGUxEjAQBgNVBAoMCXRlcnJhZm9ybTEVMBMGA1UEAwwMdGVycmFmb3JtLmlvMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA3L9L5szT4+FLykTFNyyPjy/k3BQTYAfRQzP2dhnsuUKm3cdPC0NyZ+wEXIUGhoDO2YG6EYChOl8fsDqDOjloSUGKqYw++nlpHIuUgJx8IxxG2XkALCjFU7EmF+w7kn76d0ezpEIYxnLP+KG2DVornoEt1aLhv1MLmpgEZZPhDbMSLhSYWeTVRMayXLwqtfgnDumQSB+8d/1JuJqrSI4pD12JozVThzb6hsjfb6RMX4epPmrGn0PbTPEEA6awmsxBCXB0s13nNQt/O0hLM2agwvAyozilQV+s616Ckgk6DJoUkqZhDy7vPYMIRSr98fBws6zkrV6tTLjmD8xAvobePQIDAQABo1AwTjAdBgNVHQ4EFgQUXIqO421zMMmbcRRX9wctZFCQuPIwHwYDVR0jBBgwFoAUXIqO421zMMmbcRRX9wctZFCQuPIwDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAr82NeT3BYJOKLlUL6Om5LjUF66ewcJjG9ltdvyQwVneMcq7t5UAPxgChzqNRVk4da8PzkXpjBJyWezHupdJNX3XqeUk2kSxqQ6/gmhqvfI3y7djrwoO6jvMEY26WqtkTNORWDP3THJJVimC3zV+KMU5UBVrEzhOVhHSU709lBP75o0BBn3xGsPqSq9k8IotIFfyAc6a+XP3+ZMpvh7wqAUml7vWa5wlcXExCx39h1balfDSLGNC4swWPCp9AMnQR0p+vMay9hNP1Eh+9QYUai14d5KS3cFV+KxE1cJR5HD/iLltnnOEbpMsB0eVOZWkFvE7Y5lW0oVSAfin5TwTJMQ=="

  metadata = <<METADATA
    {
        "foo": "bar2"
    }
METADATA

  key_vault_key {
    key_name     = azurerm_key_vault_key.test2.name
    key_vault_id = azurerm_key_vault.test2.id
    key_version  = azurerm_key_vault_key.test2.version
  }
}
`, template, data.RandomString, data.RandomString, data.RandomInteger)
}
