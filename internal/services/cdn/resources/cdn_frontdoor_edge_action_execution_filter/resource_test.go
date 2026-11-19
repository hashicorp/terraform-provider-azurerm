// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_edge_action_execution_filter_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactionexecutionfilters"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CdnFrontdoorEdgeActionExecutionFilterResource struct{}

func TestAccCdnFrontdoorEdgeActionExecutionFilter_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_edge_action_execution_filter", "test")
	r := CdnFrontdoorEdgeActionExecutionFilterResource{}

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

func TestAccCdnFrontdoorEdgeActionExecutionFilter_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_edge_action_execution_filter", "test")
	r := CdnFrontdoorEdgeActionExecutionFilterResource{}

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

func TestAccCdnFrontdoorEdgeActionExecutionFilter_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_edge_action_execution_filter", "test")
	r := CdnFrontdoorEdgeActionExecutionFilterResource{}

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

func TestAccCdnFrontdoorEdgeActionExecutionFilter_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_edge_action_execution_filter", "test")
	r := CdnFrontdoorEdgeActionExecutionFilterResource{}

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

func TestAccCdnFrontdoorEdgeActionExecutionFilter_versionFromAnotherEdgeAction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_edge_action_execution_filter", "test")
	r := CdnFrontdoorEdgeActionExecutionFilterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.versionFromAnotherEdgeAction(data),
			ExpectError: regexp.MustCompile("must be a version of the Edge Action specified in `cdn_frontdoor_edge_action_id`"),
		},
	})
}

func TestAccCdnFrontdoorEdgeActionExecutionFilter_replaceReferencedVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_edge_action_execution_filter", "test")
	r := CdnFrontdoorEdgeActionExecutionFilterResource{}

	// renaming the referenced version replaces it, which the API only allows once the execution filter no longer
	// references it - `replace_triggered_by` recreates the execution filter together with the version
	data.ResourceTestIgnoreRecreate(t, r, []acceptance.TestStep{
		{
			Config: r.replaceTriggeredByVersion(data, "canary"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.replaceTriggeredByVersion(data, "canaryrenamed"),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction("azurerm_cdn_frontdoor_edge_action_version.canary", plancheck.ResourceActionDestroyBeforeCreate),
					plancheck.ExpectResourceAction(data.ResourceName, plancheck.ResourceActionDestroyBeforeCreate),
				},
			},
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (CdnFrontdoorEdgeActionExecutionFilterResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := edgeactionexecutionfilters.ParseExecutionFilterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Cdn.EdgeActionFiltersClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (CdnFrontdoorEdgeActionExecutionFilterResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-eaf-%[1]d"
  location = "%[2]s"
}

resource "azurerm_cdn_frontdoor_edge_action" "test" {
  name                = "acctestea%[1]d"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_cdn_frontdoor_edge_action_version" "default" {
  name                         = "default%[1]d"
  cdn_frontdoor_edge_action_id = azurerm_cdn_frontdoor_edge_action.test.id
  deployment_type              = "file"
}

resource "azurerm_cdn_frontdoor_edge_action_version" "first" {
  name                         = "first%[1]d"
  cdn_frontdoor_edge_action_id = azurerm_cdn_frontdoor_edge_action.test.id
  deployment_type              = "file"

  # the first version of an Edge Action becomes its default version, which can only be
  # deleted once the other versions are gone, so this version is created after it
  depends_on = [azurerm_cdn_frontdoor_edge_action_version.default]
}

resource "azurerm_cdn_frontdoor_edge_action_version" "second" {
  name                         = "second%[1]d"
  cdn_frontdoor_edge_action_id = azurerm_cdn_frontdoor_edge_action.test.id
  deployment_type              = "file"

  depends_on = [azurerm_cdn_frontdoor_edge_action_version.default]
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r CdnFrontdoorEdgeActionExecutionFilterResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_edge_action_execution_filter" "test" {
  name                                 = "acctesteaf%d"
  cdn_frontdoor_edge_action_id         = azurerm_cdn_frontdoor_edge_action.test.id
  cdn_frontdoor_edge_action_version_id = azurerm_cdn_frontdoor_edge_action_version.first.id
  header_name                          = "x-canary"
  header_value                         = "true"
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontdoorEdgeActionExecutionFilterResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_edge_action_execution_filter" "import" {
  name                                 = azurerm_cdn_frontdoor_edge_action_execution_filter.test.name
  cdn_frontdoor_edge_action_id         = azurerm_cdn_frontdoor_edge_action_execution_filter.test.cdn_frontdoor_edge_action_id
  cdn_frontdoor_edge_action_version_id = azurerm_cdn_frontdoor_edge_action_execution_filter.test.cdn_frontdoor_edge_action_version_id
  header_name                          = azurerm_cdn_frontdoor_edge_action_execution_filter.test.header_name
  header_value                         = azurerm_cdn_frontdoor_edge_action_execution_filter.test.header_value
}
`, r.basic(data))
}

func (r CdnFrontdoorEdgeActionExecutionFilterResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_edge_action_execution_filter" "test" {
  name                                 = "acctesteaf%d"
  cdn_frontdoor_edge_action_id         = azurerm_cdn_frontdoor_edge_action.test.id
  cdn_frontdoor_edge_action_version_id = azurerm_cdn_frontdoor_edge_action_version.second.id
  header_name                          = "x-test"
  header_value                         = "execution"

  tags = {
    Environment = "Test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontdoorEdgeActionExecutionFilterResource) versionFromAnotherEdgeAction(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cdn_frontdoor_edge_action" "other" {
  name                = "acctesteaother%[2]d"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_cdn_frontdoor_edge_action_execution_filter" "test" {
  name                                 = "acctesteaf%[2]d"
  cdn_frontdoor_edge_action_id         = azurerm_cdn_frontdoor_edge_action.other.id
  cdn_frontdoor_edge_action_version_id = azurerm_cdn_frontdoor_edge_action_version.first.id
  header_name                          = "x-canary"
  header_value                         = "true"
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontdoorEdgeActionExecutionFilterResource) replaceTriggeredByVersion(data acceptance.TestData, versionName string) string {
	// an Edge Action can have at most 3 versions, so this test uses its own Edge Action with only 2 versions
	// rather than the template, leaving room for the replacement version which is created before the old one is deleted
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-eaf-%[1]d"
  location = "%[2]s"
}

resource "azurerm_cdn_frontdoor_edge_action" "test" {
  name                = "acctestea%[1]d"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_cdn_frontdoor_edge_action_version" "default" {
  name                         = "default%[1]d"
  cdn_frontdoor_edge_action_id = azurerm_cdn_frontdoor_edge_action.test.id
  deployment_type              = "file"
}

resource "azurerm_cdn_frontdoor_edge_action_version" "canary" {
  name                         = "%[3]s%[1]d"
  cdn_frontdoor_edge_action_id = azurerm_cdn_frontdoor_edge_action.test.id
  deployment_type              = "file"

  depends_on = [azurerm_cdn_frontdoor_edge_action_version.default]
}

resource "azurerm_cdn_frontdoor_edge_action_execution_filter" "test" {
  name                                 = "acctesteaf%[1]d"
  cdn_frontdoor_edge_action_id         = azurerm_cdn_frontdoor_edge_action.test.id
  cdn_frontdoor_edge_action_version_id = azurerm_cdn_frontdoor_edge_action_version.canary.id
  header_name                          = "x-canary"
  header_value                         = "true"

  lifecycle {
    replace_triggered_by = [azurerm_cdn_frontdoor_edge_action_version.canary.id]
  }
}
`, data.RandomInteger, data.Locations.Primary, versionName)
}
