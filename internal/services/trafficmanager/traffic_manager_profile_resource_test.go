// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package trafficmanager_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type TrafficManagerProfileResource struct{}

func TestAccTrafficManagerProfile_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_profile", "test")
	r := TrafficManagerProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "Geographic"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				acceptance.TestCheckResourceAttr(data.ResourceName, "traffic_routing_method", "Geographic"),
				acceptance.TestCheckResourceAttr(data.ResourceName, "fqdn", fmt.Sprintf("acctest-tmp-%d.trafficmanager.net", data.RandomInteger)),
			),
		},
		data.ImportStep(),
	})
}

func TestAccTrafficManagerProfile_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_profile", "test")
	r := TrafficManagerProfileResource{}

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

func TestAccTrafficManagerProfile_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_profile", "test")
	r := TrafficManagerProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccTrafficManagerProfile_updateEnsureDoNotEraseEndpoints(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_profile", "test")
	r := TrafficManagerProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeWithEndpoint(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeUpdatedWithEndpoint(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccTrafficManagerProfile_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_profile", "test")
	r := TrafficManagerProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "Geographic"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("traffic_routing_method").HasValue("Geographic"),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccTrafficManagerProfile_cycleMethod(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_profile", "test")
	r := TrafficManagerProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "Geographic"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("traffic_routing_method").HasValue("Geographic"),
				check.That(data.ResourceName).Key("fqdn").HasValue(fmt.Sprintf("acctest-tmp-%d.trafficmanager.net", data.RandomInteger)),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, "Weighted"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("traffic_routing_method").HasValue("Weighted"),
				check.That(data.ResourceName).Key("fqdn").HasValue(fmt.Sprintf("acctest-tmp-%d.trafficmanager.net", data.RandomInteger)),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, "Subnet"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("traffic_routing_method").HasValue("Subnet"),
				check.That(data.ResourceName).Key("fqdn").HasValue(fmt.Sprintf("acctest-tmp-%d.trafficmanager.net", data.RandomInteger)),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, "Priority"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("traffic_routing_method").HasValue("Priority"),
				check.That(data.ResourceName).Key("fqdn").HasValue(fmt.Sprintf("acctest-tmp-%d.trafficmanager.net", data.RandomInteger)),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, "Performance"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("traffic_routing_method").HasValue("Performance"),
				check.That(data.ResourceName).Key("fqdn").HasValue(fmt.Sprintf("acctest-tmp-%d.trafficmanager.net", data.RandomInteger)),
			),
		},
		data.ImportStep(),
		{
			Config: r.multiValue(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("traffic_routing_method").HasValue("MultiValue"),
				check.That(data.ResourceName).Key("fqdn").HasValue(fmt.Sprintf("acctest-tmp-%d.trafficmanager.net", data.RandomInteger)),
			),
		},
		data.ImportStep(),
	})
}

func TestAccTrafficManagerProfile_fastEndpointFailoverSettingsError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_profile", "test")
	r := TrafficManagerProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.failoverError(data),
			ExpectError: regexp.MustCompile("`timeout_in_seconds` must be between `5` and `9` when `interval_in_seconds` is set to `10`"),
		},
	})
}

func TestAccTrafficManagerProfile_fastMaxReturnSettingError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_profile", "test")
	r := TrafficManagerProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.maxReturnError(data),
			ExpectError: regexp.MustCompile("`max_return` must be specified when `traffic_routing_method` is set to `MultiValue`"),
		},
	})
}

func TestAccTrafficManagerProfile_trafficView(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_profile", "test")
	r := TrafficManagerProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withTrafficView(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("traffic_view_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.withTrafficView(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("traffic_view_enabled").HasValue("true"),
			),
		},
	})
}

func TestAccTrafficManagerProfile_updateTTL(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_traffic_manager_profile", "test")
	r := TrafficManagerProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withTTL(data, "Geographic", 0),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withTTL(data, "Geographic", 2147483647),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r TrafficManagerProfileResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := profiles.ParseTrafficManagerProfileID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.TrafficManager.ProfilesClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(true), nil
}

func (r TrafficManagerProfileResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-traffic-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r TrafficManagerProfileResource) basic(data acceptance.TestData, method string) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctest-TMP-%d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "%s"

  dns_config {
    relative_name = "acctest-tmp-%d"
    ttl           = 30
  }

  monitor_config {
    protocol = "HTTPS"
    port     = 443
    path     = "/"
  }
}
`, template, data.RandomInteger, method, data.RandomInteger)
}

func (r TrafficManagerProfileResource) multiValue(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctest-TMP-%d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "MultiValue"
  max_return             = 8

  dns_config {
    relative_name = "acctest-tmp-%d"
    ttl           = 30
  }

  monitor_config {
    protocol = "HTTPS"
    port     = 443
    path     = "/"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (r TrafficManagerProfileResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data, "Geographic")
	return fmt.Sprintf(`
%s

resource "azurerm_traffic_manager_profile" "import" {
  name                   = azurerm_traffic_manager_profile.test.name
  resource_group_name    = azurerm_traffic_manager_profile.test.resource_group_name
  traffic_routing_method = azurerm_traffic_manager_profile.test.traffic_routing_method

  dns_config {
    relative_name = "acctest-tmp-%d"
    ttl           = 30
  }

  monitor_config {
    protocol = "HTTPS"
    port     = 443
    path     = "/"
  }
}
`, template, data.RandomInteger)
}

func (r TrafficManagerProfileResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctest-TMP-%d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Performance"

  dns_config {
    relative_name = "acctest-tmp-%d"
    ttl           = 30
  }

  monitor_config {
    expected_status_code_ranges = [
      "100-101",
      "301-303",
    ]

    custom_header {
      name  = "foo"
      value = "bar"
    }

    protocol = "TCP"
    port     = 777

    interval_in_seconds          = 30
    timeout_in_seconds           = 9
    tolerated_number_of_failures = 6
  }

  tags = {
    Environment = "Production"
    cost_center = "acctest"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (r TrafficManagerProfileResource) completeUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctest-TMP-%d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Priority"

  dns_config {
    relative_name = "acctest-tmp-%d"
    ttl           = 30
  }

  monitor_config {
    expected_status_code_ranges = [
      "302-304",
    ]

    custom_header {
      name  = "foo2"
      value = "bar2"
    }

    protocol = "HTTPS"
    port     = 442
    path     = "/"

    interval_in_seconds          = 30
    timeout_in_seconds           = 6
    tolerated_number_of_failures = 3
  }

  tags = {
    Environment = "staging"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (r TrafficManagerProfileResource) endpointResource(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_traffic_manager_external_endpoint" "test" {
  name       = "acctestend-external%d"
  profile_id = azurerm_traffic_manager_profile.test.id
  target     = "pluginsdk.io"
  weight     = 100
}
`, data.RandomInteger)
}

func (r TrafficManagerProfileResource) completeWithEndpoint(data acceptance.TestData) string {
	template := r.template(data)
	endpoint := r.endpointResource(data)
	return fmt.Sprintf(`
%s

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctest-TMP-%d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Weighted"

  dns_config {
    relative_name = "acctest-tmp-%d"
    ttl           = 30
  }

  monitor_config {
    expected_status_code_ranges = [
      "100-101",
      "301-303",
    ]

    custom_header {
      name  = "foo"
      value = "bar"
    }

    protocol = "TCP"
    port     = 777

    interval_in_seconds          = 30
    timeout_in_seconds           = 9
    tolerated_number_of_failures = 6
  }

  tags = {
    Environment = "Production"
    cost_center = "acctest"
  }
}

%s
`, template, data.RandomInteger, data.RandomInteger, endpoint)
}

func (r TrafficManagerProfileResource) completeUpdatedWithEndpoint(data acceptance.TestData) string {
	template := r.template(data)
	endpoint := r.endpointResource(data)
	return fmt.Sprintf(`
%s

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctest-TMP-%d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Priority"

  dns_config {
    relative_name = "acctest-tmp-%d"
    ttl           = 30
  }

  monitor_config {
    expected_status_code_ranges = [
      "302-304",
    ]

    custom_header {
      name  = "foo2"
      value = "bar2"
    }

    protocol = "HTTPS"
    port     = 442
    path     = "/"

    interval_in_seconds          = 30
    timeout_in_seconds           = 6
    tolerated_number_of_failures = 3
  }

  tags = {
    Environment = "staging"
  }
}

%s
`, template, data.RandomInteger, data.RandomInteger, endpoint)
}

func (r TrafficManagerProfileResource) failoverError(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctest-TMP-%d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Performance"

  dns_config {
    relative_name = "acctest-tmp-%d"
    ttl           = 30
  }

  monitor_config {
    protocol                     = "HTTPS"
    port                         = 443
    path                         = "/"
    interval_in_seconds          = 10
    timeout_in_seconds           = 10
    tolerated_number_of_failures = 3
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (r TrafficManagerProfileResource) maxReturnError(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctest-TMP-%d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "MultiValue"

  dns_config {
    relative_name = "acctest-tmp-%d"
    ttl           = 30
  }

  monitor_config {
    protocol                     = "HTTPS"
    port                         = 443
    path                         = "/"
    interval_in_seconds          = 10
    timeout_in_seconds           = 8
    tolerated_number_of_failures = 3
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (r TrafficManagerProfileResource) withTrafficView(data acceptance.TestData, enabled bool) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctest-TMP-%d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "Geographic"

  dns_config {
    relative_name = "acctest-tmp-%d"
    ttl           = 30
  }

  monitor_config {
    protocol = "HTTPS"
    port     = 443
    path     = "/"
  }
  traffic_view_enabled = %t
}
`, template, data.RandomInteger, data.RandomInteger, enabled)
}

func (r TrafficManagerProfileResource) withTTL(data acceptance.TestData, method string, ttl int) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "acctest-TMP-%d"
  resource_group_name    = azurerm_resource_group.test.name
  traffic_routing_method = "%s"

  dns_config {
    relative_name = "acctest-tmp-%d"
    ttl           = %d
  }

  monitor_config {
    protocol = "HTTPS"
    port     = 443
    path     = "/"
  }
}
`, template, data.RandomInteger, method, data.RandomInteger, ttl)
}
