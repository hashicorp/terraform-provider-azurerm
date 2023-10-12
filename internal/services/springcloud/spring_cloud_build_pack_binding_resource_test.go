// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package springcloud_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SpringCloudBuildPackBindingResource struct{}

func TestAccSpringCloudBuildPackBinding_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_build_pack_binding", "test")
	r := SpringCloudBuildPackBindingResource{}
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

func TestAccSpringCloudBuildPackBinding_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_build_pack_binding", "test")
	r := SpringCloudBuildPackBindingResource{}
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

func TestAccSpringCloudBuildPackBinding_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_build_pack_binding", "test")
	r := SpringCloudBuildPackBindingResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("launch.0.secrets.%", "launch.0.secrets.connection-string"),
	})
}

func TestAccSpringCloudBuildPackBinding_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_build_pack_binding", "test")
	r := SpringCloudBuildPackBindingResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("launch.0.secrets.%", "launch.0.secrets.connection-string"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r SpringCloudBuildPackBindingResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SpringCloudBuildPackBindingID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.AppPlatform.BuildPackBindingClient.Get(ctx, id.ResourceGroup, id.SpringName, id.BuildServiceName, id.BuilderName, id.BuildPackBindingName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r SpringCloudBuildPackBindingResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-spring-%[2]d"
  location = "%[1]s"
}

resource "azurerm_spring_cloud_service" "test" {
  name                = "acctest-sc-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "E0"
}

resource "azurerm_spring_cloud_builder" "test" {
  name                    = "acc-%[3]s"
  spring_cloud_service_id = azurerm_spring_cloud_service.test.id
  build_pack_group {
    name           = "mix"
    build_pack_ids = ["tanzu-buildpacks/java-azure"]
  }

  stack {
    id      = "io.buildpacks.stacks.bionic"
    version = "base"
  }
}
`, data.Locations.Primary, data.RandomInteger, data.RandomStringOfLength(5))
}

func (r SpringCloudBuildPackBindingResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_build_pack_binding" "test" {
  name                    = "acc-%s"
  spring_cloud_builder_id = azurerm_spring_cloud_builder.test.id
  binding_type            = "ApplicationInsights"
}
`, template, data.RandomStringOfLength(5))
}

func (r SpringCloudBuildPackBindingResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_build_pack_binding" "import" {
  name                    = azurerm_spring_cloud_build_pack_binding.test.name
  spring_cloud_builder_id = azurerm_spring_cloud_build_pack_binding.test.spring_cloud_builder_id
  binding_type            = azurerm_spring_cloud_build_pack_binding.test.binding_type
}
`, config)
}

func (r SpringCloudBuildPackBindingResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_build_pack_binding" "test" {
  name                    = "acc-%s"
  spring_cloud_builder_id = azurerm_spring_cloud_builder.test.id
  binding_type            = "ApplicationInsights"
  launch {
    properties = {
      abc           = "def"
      any-string    = "any-string"
      sampling-rate = "12.0"
    }

    secrets = {
      connection-string = "XXXXXXXXXXXXXXXXX=XXXXXXXXXXXXX-XXXXXXXXXXXXXXXXXXX;XXXXXXXXXXXXXXXXX=XXXXXXXXXXXXXXXXXXX"
    }
  }
}
`, template, data.RandomStringOfLength(5))
}
