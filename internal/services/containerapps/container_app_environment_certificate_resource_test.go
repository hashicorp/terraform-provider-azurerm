// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containerapps_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/certificates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ContainerAppEnvironmentCertificateResource struct{}

func TestAccContainerAppEnvironmentCertificate_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_environment_certificate", "test")
	r := ContainerAppEnvironmentCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("certificate_blob_base64", "certificate_password"),
	})
}

func TestAccContainerAppEnvironmentCertificate_basicUpdateTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_app_environment_certificate", "test")
	r := ContainerAppEnvironmentCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("certificate_blob_base64", "certificate_password"),
		{
			Config: r.basicAddTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("certificate_blob_base64", "certificate_password"),
	})
}

func (r ContainerAppEnvironmentCertificateResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := certificates.ParseCertificateID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.ContainerApps.CertificatesClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ContainerAppEnvironmentCertificateResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_container_app_environment_certificate" "test" {
  name                         = "acctest-cacert%[2]d"
  container_app_environment_id = azurerm_container_app_environment.test.id
  certificate_blob_base64      = filebase64("testdata/testacc.pfx")
  certificate_password         = "TestAcc"
}
`, r.template(data), data.RandomInteger)
}

func (r ContainerAppEnvironmentCertificateResource) basicAddTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_container_app_environment_certificate" "test" {
  name                         = "acctest-cacert%[2]d"
  container_app_environment_id = azurerm_container_app_environment.test.id
  certificate_blob_base64      = filebase64("testdata/testacc.pfx")
  certificate_password         = "TestAcc"

  tags = {
    env = "testAcc"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ContainerAppEnvironmentCertificateResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-CAEnv-%[1]d"
  location = "%[2]s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestCAEnv-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_container_app_environment" "test" {
  name                       = "acctest-CAEnv%[1]d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}
