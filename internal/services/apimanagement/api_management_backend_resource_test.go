// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/backend"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementBackendResource struct{}

func TestAccApiManagementBackend_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_backend", "test")
	r := ApiManagementBackendResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "basic"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("protocol").HasValue("http"),
				check.That(data.ResourceName).Key("url").HasValue("https://acctest"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementBackend_circuitBreakerRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_backend", "test")
	r := ApiManagementBackendResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.circuitBreakerRuleComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.circuitBreakerRuleCompleteUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.circuitBreakerRuleBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.circuitBreakerRuleComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, "circuitbreaker"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementBackend_allProperties(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_backend", "test")
	r := ApiManagementBackendResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.allProperties(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("protocol").HasValue("http"),
				check.That(data.ResourceName).Key("url").HasValue("https://acctest"),
				check.That(data.ResourceName).Key("description").HasValue("description"),
				check.That(data.ResourceName).Key("resource_id").HasValue("https://resourceid"),
				check.That(data.ResourceName).Key("title").HasValue("title"),
				check.That(data.ResourceName).Key("credentials.#").HasValue("1"),
				check.That(data.ResourceName).Key("credentials.0.authorization.0.parameter").HasValue("parameter"),
				check.That(data.ResourceName).Key("credentials.0.authorization.0.scheme").HasValue("scheme"),
				check.That(data.ResourceName).Key("credentials.0.certificate.0").Exists(),
				check.That(data.ResourceName).Key("credentials.0.header.header1").HasValue("header1value1,header1value2"),
				check.That(data.ResourceName).Key("credentials.0.header.header2").HasValue("header2value1,header2value2"),
				check.That(data.ResourceName).Key("credentials.0.query.query1").HasValue("query1value1,query1value2"),
				check.That(data.ResourceName).Key("credentials.0.query.query2").HasValue("query2value1,query2value2"),
				check.That(data.ResourceName).Key("proxy.#").HasValue("1"),
				check.That(data.ResourceName).Key("proxy.0.url").HasValue("http://192.168.1.1:8080"),
				check.That(data.ResourceName).Key("proxy.0.username").HasValue("username"),
				check.That(data.ResourceName).Key("proxy.0.password").HasValue("password"),
				check.That(data.ResourceName).Key("tls.#").HasValue("1"),
				check.That(data.ResourceName).Key("tls.0.validate_certificate_chain").HasValue("false"),
				check.That(data.ResourceName).Key("tls.0.validate_certificate_name").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementBackend_credentialsNoCertificate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_backend", "test")
	r := ApiManagementBackendResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.credentialsNoCertificate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementBackend_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_backend", "test")
	r := ApiManagementBackendResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "update"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("protocol").HasValue("http"),
				check.That(data.ResourceName).Key("url").HasValue("https://acctest"),
			),
		},
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("protocol").HasValue("soap"),
				check.That(data.ResourceName).Key("url").HasValue("https://updatedacctest"),
				check.That(data.ResourceName).Key("description").HasValue("description"),
				check.That(data.ResourceName).Key("resource_id").HasValue("https://resourceid"),
				check.That(data.ResourceName).Key("proxy.#").HasValue("1"),
				check.That(data.ResourceName).Key("proxy.0.url").HasValue("http://192.168.1.1:8080"),
				check.That(data.ResourceName).Key("proxy.0.username").HasValue("username"),
				check.That(data.ResourceName).Key("proxy.0.password").HasValue("password"),
				check.That(data.ResourceName).Key("tls.#").HasValue("1"),
				check.That(data.ResourceName).Key("tls.0.validate_certificate_chain").HasValue("false"),
				check.That(data.ResourceName).Key("tls.0.validate_certificate_name").HasValue("true"),
			),
		},
		{
			Config: r.basic(data, "update"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("protocol").HasValue("http"),
				check.That(data.ResourceName).Key("url").HasValue("https://acctest"),
				check.That(data.ResourceName).Key("description").HasValue(""),
				check.That(data.ResourceName).Key("resource_id").HasValue(""),
				check.That(data.ResourceName).Key("proxy.#").HasValue("0"),
				check.That(data.ResourceName).Key("tls.#").HasValue("0"),
			),
		},
	})
}

func TestAccApiManagementBackend_serviceFabric(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_backend", "test")
	r := ApiManagementBackendResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.serviceFabric(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("service_fabric_cluster.0.client_certificate_thumbprint").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementBackend_serviceFabricCluster(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_backend", "test")
	r := ApiManagementBackendResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.serviceFabricCluster(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementBackend_serviceFabricClientCertificateId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_backend", "test")
	r := ApiManagementBackendResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.serviceFabricClientCertificateId(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementBackend_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_backend", "test")
	r := ApiManagementBackendResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config: func(d acceptance.TestData) string {
				return r.basic(d, "disappears")
			},
			TestResource: r,
		}),
	})
}

func TestAccApiManagementBackend_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_backend", "test")
	r := ApiManagementBackendResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "import"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (ApiManagementBackendResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := backend.ParseBackendID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.BackendClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil && resp.Model.Id != nil), nil
}

func (r ApiManagementBackendResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := backend.ParseBackendID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.ApiManagement.BackendClient.Delete(ctx, *id, backend.DeleteOperationOptions{})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(true), nil
		}
		return nil, fmt.Errorf("deleting Backend: %+v", err)
	}

	return pointer.To(true), nil
}

func (r ApiManagementBackendResource) basic(data acceptance.TestData, testName string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_backend" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  protocol            = "http"
  url                 = "https://acctest"
}
`, r.template(data, testName), data.RandomInteger)
}

func (r ApiManagementBackendResource) basicForResourceIdentity(data acceptance.TestData) string {
	return r.basic(data, "basic")
}

func (r ApiManagementBackendResource) circuitBreakerRuleBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_api_management_backend" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  protocol            = "http"
  url                 = "https://acctest"
  description         = "Test backend with circuit breaker rule"

  circuit_breaker_rule {
    name          = "test-circuit-breaker"
    trip_duration = "PT30S"
    failure_condition {
      count             = 5
      interval_duration = "PT1M"
      status_code_range {
        min = 200
        max = 299
      }
    }
  }
}
`, r.template(data, "circuitbreaker"), data.RandomInteger)
}

func (r ApiManagementBackendResource) circuitBreakerRuleComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_api_management_backend" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  protocol            = "http"
  url                 = "https://acctest"
  description         = "Test backend with circuit breaker rule"

  circuit_breaker_rule {
    name                       = "test-circuit-breaker"
    trip_duration              = "PT30S"
    accept_retry_after_enabled = true
    failure_condition {
      count             = 5
      interval_duration = "PT1M"
      status_code_range {
        min = 200
        max = 299
      }
      status_code_range {
        min = 300
        max = 399
      }
      error_reasons = ["SubscriptionKeyInvalid", "ClientConnectionFailure", "OperationNotFound"]
    }
  }
}
`, r.template(data, "circuitbreaker"), data.RandomInteger)
}

func (r ApiManagementBackendResource) circuitBreakerRuleCompleteUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_api_management_backend" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  protocol            = "http"
  url                 = "https://acctest"
  description         = "Test backend with circuit breaker rule"

  circuit_breaker_rule {
    name                       = "test-circuit-breaker-update"
    trip_duration              = "PT40S"
    accept_retry_after_enabled = false
    failure_condition {
      percentage        = 2
      interval_duration = "PT2M"
      status_code_range {
        min = 300
        max = 399
      }
      error_reasons = ["BackendConnectionFailure"]
    }
  }
}
`, r.template(data, "circuitbreaker"), data.RandomInteger)
}

func (r ApiManagementBackendResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_backend" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  protocol            = "soap"
  url                 = "https://updatedacctest"
  description         = "description"
  resource_id         = "https://resourceid"
  proxy {
    url      = "http://192.168.1.1:8080"
    username = "username"
    password = "password"
  }
  tls {
    validate_certificate_chain = false
    validate_certificate_name  = true
  }
}
`, r.template(data, "update"), data.RandomInteger)
}

func (r ApiManagementBackendResource) allProperties(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_certificate" "test" {
  name                = "example-cert"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  data                = filebase64("testdata/keyvaultcert.pfx")
  password            = ""
}

resource "azurerm_api_management_backend" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  protocol            = "http"
  url                 = "https://acctest"
  description         = "description"
  resource_id         = "https://resourceid"
  title               = "title"
  credentials {
    authorization {
      parameter = "parameter"
      scheme    = "scheme"
    }
    certificate = [
      azurerm_api_management_certificate.test.thumbprint,
    ]
    header = {
      header1 = "header1value1,header1value2"
      header2 = "header2value1,header2value2"
    }
    query = {
      query1 = "query1value1,query1value2"
      query2 = "query2value1,query2value2"
    }
  }
  proxy {
    url      = "http://192.168.1.1:8080"
    username = "username"
    password = "password"
  }
  tls {
    validate_certificate_chain = false
    validate_certificate_name  = true
  }
}
`, r.template(data, "all"), data.RandomInteger)
}

func (r ApiManagementBackendResource) serviceFabric(data acceptance.TestData) string {
	// nolint: dupword
	return fmt.Sprintf(` 
%s

resource "azurerm_api_management_certificate" "test" {
  name                = "example-cert"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  data                = filebase64("testdata/keyvaultcert.pfx")
  password            = ""
}

resource "azurerm_api_management_backend" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  protocol            = "http"
  url                 = "fabric:/mytestapp/acctest"
  service_fabric_cluster {
    client_certificate_thumbprint = azurerm_api_management_certificate.test.thumbprint
    management_endpoints = [
      "https://acctestsf.com",
    ]
    max_partition_resolution_retries = 5
    server_certificate_thumbprints = [
      azurerm_api_management_certificate.test.thumbprint,
      azurerm_api_management_certificate.test.thumbprint,
    ]
  }
}
`, r.template(data, "sf"), data.RandomInteger)
}

func (r ApiManagementBackendResource) serviceFabricCluster(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_certificate" "test" {
  name                = "acctest-cert-%[2]d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  data                = filebase64("testdata/keyvaultcert.pfx")
  password            = ""
}

resource "azurerm_api_management_backend" "test" {
  name                = "acctestapi-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  protocol            = "http"
  url                 = "fabric:/mytestapp/acctest"
  service_fabric_cluster {
    client_certificate_thumbprint = azurerm_api_management_certificate.test.thumbprint
    server_x509_name {
      name                          = "test"
      issuer_certificate_thumbprint = azurerm_api_management_certificate.test.thumbprint
    }
    management_endpoints = [
      "https://acctestsf.com",
    ]
    max_partition_resolution_retries = 5
  }
}
`, r.template(data, "sf"), data.RandomInteger)
}

func (r ApiManagementBackendResource) serviceFabricClientCertificateId(data acceptance.TestData) string {
	// nolint: dupword
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_certificate" "test" {
  name                = "example-cert"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  data                = filebase64("testdata/keyvaultcert.pfx")
  password            = ""
}

resource "azurerm_api_management_backend" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  protocol            = "http"
  url                 = "fabric:/mytestapp/acctest"
  service_fabric_cluster {
    client_certificate_id = azurerm_api_management_certificate.test.id
    management_endpoints = [
      "https://acctestsf.com",
    ]
    max_partition_resolution_retries = 5
    server_certificate_thumbprints = [
      azurerm_api_management_certificate.test.thumbprint,
      azurerm_api_management_certificate.test.thumbprint,
    ]
  }
}
`, r.template(data, "sf"), data.RandomInteger)
}

func (r ApiManagementBackendResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_backend" "import" {
  name                = azurerm_api_management_backend.test.name
  resource_group_name = azurerm_api_management_backend.test.resource_group_name
  api_management_name = azurerm_api_management_backend.test.api_management_name
  protocol            = azurerm_api_management_backend.test.protocol
  url                 = azurerm_api_management_backend.test.url
}
`, r.basic(data, "import"))
}

func (ApiManagementBackendResource) template(data acceptance.TestData, testName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d-%s"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Consumption_0"
}
`, data.RandomInteger, testName, data.Locations.Primary, data.RandomInteger, testName)
}

func (r ApiManagementBackendResource) credentialsNoCertificate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_backend" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  protocol            = "http"
  url                 = "https://acctest"
  description         = "description"
  resource_id         = "https://resourceid"
  title               = "title"
  credentials {
    authorization {
      parameter = "parameter"
      scheme    = "scheme"
    }
    header = {
      header1 = "header1value1,header1value2"
      header2 = "header2value1,header2value2"
    }
    query = {
      query1 = "query1value1,query1value2"
      query2 = "query2value1,query2value2"
    }
  }
  proxy {
    url      = "http://192.168.1.1:8080"
    username = "username"
    password = "password"
  }
  tls {
    validate_certificate_chain = false
    validate_certificate_name  = true
  }
}
`, r.template(data, "all"), data.RandomInteger)
}
