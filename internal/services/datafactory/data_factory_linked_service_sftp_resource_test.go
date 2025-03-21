// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LinkedServiceSFTPResource struct{}

func TestAccDataFactoryLinkedServiceSFTP_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_service_sftp", "test")
	r := LinkedServiceSFTPResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
	})
}

func TestAccDataFactoryLinkedServiceSFTP_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_service_sftp", "test")
	r := LinkedServiceSFTPResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
	})
}

func TestAccDataFactoryLinkedServiceSFTP_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_service_sftp", "test")
	r := LinkedServiceSFTPResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.update1(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("parameters.%").HasValue("2"),
				check.That(data.ResourceName).Key("annotations.#").HasValue("3"),
				check.That(data.ResourceName).Key("additional_properties.%").HasValue("2"),
				check.That(data.ResourceName).Key("description").HasValue("test description"),
			),
		},
		data.ImportStep("password"),
		{
			Config: r.update2(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("parameters.%").HasValue("3"),
				check.That(data.ResourceName).Key("annotations.#").HasValue("2"),
				check.That(data.ResourceName).Key("additional_properties.%").HasValue("1"),
				check.That(data.ResourceName).Key("description").HasValue("test description 2"),
			),
		},
		data.ImportStep("password"),
	})
}

func (t LinkedServiceSFTPResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.LinkedServiceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DataFactory.LinkedServiceClient.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("reading Data Factory SFTP (%s): %+v", *id, err)
	}

	return pointer.To(resp.ID != nil), nil
}

func (t LinkedServiceSFTPResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_data_factory_linked_service_sftp" "test" {
  name                = "acctestlssftp%d"
  data_factory_id     = azurerm_data_factory.test.id
  authentication_type = "Basic"
  host                = "http://www.bing.com"
  port                = 22
  username            = "foo"
  password            = "bar"
}
`, t.template(data), data.RandomInteger)
}

func (t LinkedServiceSFTPResource) update1(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_data_factory_linked_service_sftp" "test" {
  name                = "acctestlssftp%d"
  data_factory_id     = azurerm_data_factory.test.id
  authentication_type = "Basic"
  host                = "http://www.bing.com"
  port                = 22
  username            = "foo"
  password            = "bar"
  annotations         = ["test1", "test2", "test3"]
  description         = "test description"

  parameters = {
    foo = "test1"
    bar = "test2"
  }

  additional_properties = {
    foo = "test1"
    bar = "test2"
  }
}
`, t.template(data), data.RandomInteger)
}

func (t LinkedServiceSFTPResource) update2(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_data_factory_linked_service_sftp" "test" {
  name                     = "acctestlssftp%d"
  data_factory_id          = azurerm_data_factory.test.id
  authentication_type      = "Basic"
  host                     = "http://www.bing.com"
  port                     = 22
  username                 = "foo"
  password                 = "bar"
  annotations              = ["test1", "test2"]
  description              = "test description 2"
  skip_host_key_validation = true
  host_key_fingerprint     = "fingerprint"

  parameters = {
    foo  = "test1"
    bar  = "test2"
    buzz = "test3"
  }

  additional_properties = {
    foo = "test1"
  }
}
`, t.template(data), data.RandomInteger)
}

func (t LinkedServiceSFTPResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_data_factory_integration_runtime_azure" "test" {
  data_factory_id = azurerm_data_factory.test.id
  location        = azurerm_resource_group.test.location
  name            = "acctestlssftp%[2]d"
}

resource "azurerm_data_factory_linked_service_sftp" "test" {
  name                     = "acctestlssftp%[2]d"
  data_factory_id          = azurerm_data_factory.test.id
  authentication_type      = "Basic"
  host                     = "http://www.bing.com"
  port                     = 22
  username                 = "foo"
  password                 = "bar"
  annotations              = ["test1", "test2"]
  description              = "test description 2"
  skip_host_key_validation = true
  host_key_fingerprint     = "fingerprint"
  integration_runtime_name = azurerm_data_factory_integration_runtime_azure.test.name

  parameters = {
    foo = "test1"
    bar = "test2"
  }

  additional_properties = {
    foo = "test1"
    bar = "test1"
  }
}
`, t.template(data), data.RandomInteger)
}

func (LinkedServiceSFTPResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%[1]d"
  location = "%[2]s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

`, data.RandomInteger, data.Locations.Primary)
}
