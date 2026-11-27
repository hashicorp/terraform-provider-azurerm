// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_edge_action_version_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactionversions"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CdnFrontdoorEdgeActionVersionResource struct{}

func TestAccCdnFrontdoorEdgeActionVersion_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_edge_action_version", "test")
	r := CdnFrontdoorEdgeActionVersionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_version_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCdnFrontdoorEdgeActionVersion_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_edge_action_version", "test")
	r := CdnFrontdoorEdgeActionVersionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, "first"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("code_content_base64", "code_file_name"),
	})
}

func TestAccCdnFrontdoorEdgeActionVersion_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_edge_action_version", "test")
	r := CdnFrontdoorEdgeActionVersionResource{}

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

func TestAccCdnFrontdoorEdgeActionVersion_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_edge_action_version", "test")
	r := CdnFrontdoorEdgeActionVersionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// deploying code to a default version which has no code yet is an in-place update
			Config: r.complete(data, "first"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("code_content_base64", "code_file_name"),
	})
}

func TestAccCdnFrontdoorEdgeActionVersion_updateDefaultVersionCode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_edge_action_version", "test")
	r := CdnFrontdoorEdgeActionVersionResource{}

	data.ResourceTestIgnoreRecreate(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, "first"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_version_enabled").HasValue("true"),
			),
		},
		data.ImportStep("code_content_base64", "code_file_name"),
		{
			// the deployed code of the default version can't be redeployed, so the version is recreated
			Config: r.complete(data, "second"),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(data.ResourceName, plancheck.ResourceActionDestroyBeforeCreate),
				},
			},
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("code_content_base64", "code_file_name"),
	})
}

func TestAccCdnFrontdoorEdgeActionVersion_updateNonDefaultVersionCode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_edge_action_version", "second")
	r := CdnFrontdoorEdgeActionVersionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.nonDefaultVersionWithCode(data, "first"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_version_enabled").HasValue("false"),
			),
		},
		data.ImportStep("code_content_base64", "code_file_name"),
		{
			// the code of a version which isn't the default can be redeployed in-place
			Config: r.nonDefaultVersionWithCode(data, "second"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("code_content_base64", "code_file_name"),
	})
}

func TestAccCdnFrontdoorEdgeActionVersion_updateTagsWithDeployedCode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_edge_action_version", "test")
	r := CdnFrontdoorEdgeActionVersionResource{}
	secondVersion := "azurerm_cdn_frontdoor_edge_action_version.second"

	// both the default version (`test`) and a version which isn't the default (`second`) have code deployed,
	// only their tags are changed and then removed, which must update both versions in-place without redeploying
	// the code. `ResourceTestIgnoreRecreate` is used since `ResourceTest` would replace these plan checks
	expectUpdateInPlace := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(data.ResourceName, plancheck.ResourceActionUpdate),
			plancheck.ExpectResourceAction(secondVersion, plancheck.ResourceActionUpdate),
		},
	}

	data.ResourceTestIgnoreRecreate(t, r, []acceptance.TestStep{
		{
			Config: r.tagsWithDeployedCode(data, ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_version_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("validation_status").HasValue("Succeeded"),
				check.That(secondVersion).ExistsInAzure(r),
				check.That(secondVersion).Key("default_version_enabled").HasValue("false"),
				check.That(secondVersion).Key("validation_status").HasValue("Succeeded"),
			),
		},
		data.ImportStep("code_content_base64", "code_file_name"),
		data.ImportStepFor(secondVersion, "code_content_base64", "code_file_name"),
		{
			Config:           r.tagsWithDeployedCode(data, "single"),
			ConfigPlanChecks: expectUpdateInPlace,
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(secondVersion).ExistsInAzure(r),
			),
		},
		data.ImportStep("code_content_base64", "code_file_name"),
		data.ImportStepFor(secondVersion, "code_content_base64", "code_file_name"),
		{
			Config:           r.tagsWithDeployedCode(data, "multiple"),
			ConfigPlanChecks: expectUpdateInPlace,
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(secondVersion).ExistsInAzure(r),
			),
		},
		data.ImportStep("code_content_base64", "code_file_name"),
		data.ImportStepFor(secondVersion, "code_content_base64", "code_file_name"),
		{
			Config:           r.tagsWithDeployedCode(data, ""),
			ConfigPlanChecks: expectUpdateInPlace,
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("validation_status").HasValue("Succeeded"),
				check.That(secondVersion).ExistsInAzure(r),
				check.That(secondVersion).Key("validation_status").HasValue("Succeeded"),
			),
		},
		data.ImportStep("code_content_base64", "code_file_name"),
		data.ImportStepFor(secondVersion, "code_content_base64", "code_file_name"),
	})
}

func (CdnFrontdoorEdgeActionVersionResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := edgeactionversions.ParseVersionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Cdn.EdgeActionVersionsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

// code returns base64 encoded JavaScript which adds a response header with the given value
func (CdnFrontdoorEdgeActionVersionResource) code(value string) string {
	js := fmt.Sprintf(`function handler(event) {
    event.response.headers["x-acctest"] = %q;
    return event;
}`, value)
	return base64.StdEncoding.EncodeToString([]byte(js))
}

func (CdnFrontdoorEdgeActionVersionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-eav-%[1]d"
  location = "%[2]s"
}

resource "azurerm_cdn_frontdoor_edge_action" "test" {
  name                = "acctestea%[1]d"
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r CdnFrontdoorEdgeActionVersionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_edge_action_version" "test" {
  name                         = "acctesteav%d"
  cdn_frontdoor_edge_action_id = azurerm_cdn_frontdoor_edge_action.test.id
  deployment_type              = "file"
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontdoorEdgeActionVersionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_edge_action_version" "import" {
  name                         = azurerm_cdn_frontdoor_edge_action_version.test.name
  cdn_frontdoor_edge_action_id = azurerm_cdn_frontdoor_edge_action_version.test.cdn_frontdoor_edge_action_id
  deployment_type              = azurerm_cdn_frontdoor_edge_action_version.test.deployment_type
}
`, r.basic(data))
}

func (r CdnFrontdoorEdgeActionVersionResource) complete(data acceptance.TestData, headerValue string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_edge_action_version" "test" {
  name                         = "acctesteav%d"
  cdn_frontdoor_edge_action_id = azurerm_cdn_frontdoor_edge_action.test.id
  deployment_type              = "file"
  code_file_name               = "edge_action.js"
  code_content_base64          = "%s"

  tags = {
    Environment = "Test"
    Header      = "%s"
  }
}
`, r.template(data), data.RandomInteger, r.code(headerValue), headerValue)
}

func (r CdnFrontdoorEdgeActionVersionResource) multipleVersions(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cdn_frontdoor_edge_action_version" "test" {
  name                         = "acctesteav%[2]d"
  cdn_frontdoor_edge_action_id = azurerm_cdn_frontdoor_edge_action.test.id
  deployment_type              = "file"
}

resource "azurerm_cdn_frontdoor_edge_action_version" "second" {
  name                         = "acctesteav2%[2]d"
  cdn_frontdoor_edge_action_id = azurerm_cdn_frontdoor_edge_action.test.id
  deployment_type              = "file"

  # the first version of an Edge Action becomes its default version, which can only be
  # deleted once the other versions are gone, so this version is created after it
  depends_on = [azurerm_cdn_frontdoor_edge_action_version.test]
}
`, r.template(data), data.RandomInteger)
}

func (r CdnFrontdoorEdgeActionVersionResource) nonDefaultVersionWithCode(data acceptance.TestData, headerValue string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cdn_frontdoor_edge_action_version" "test" {
  name                         = "acctesteav%[2]d"
  cdn_frontdoor_edge_action_id = azurerm_cdn_frontdoor_edge_action.test.id
  deployment_type              = "file"
}

resource "azurerm_cdn_frontdoor_edge_action_version" "second" {
  name                         = "acctesteav2%[2]d"
  cdn_frontdoor_edge_action_id = azurerm_cdn_frontdoor_edge_action.test.id
  deployment_type              = "file"
  code_file_name               = "edge_action.js"
  code_content_base64          = "%[3]s"

  depends_on = [azurerm_cdn_frontdoor_edge_action_version.test]
}
`, r.template(data), data.RandomInteger, r.code(headerValue))
}

// tagsWithDeployedCode returns a default version and a version which isn't the default, both with code deployed, using
// the given set of tags: "" for no tags, "single" for one tag, or "multiple" for an updated tag plus an additional one
func (r CdnFrontdoorEdgeActionVersionResource) tagsWithDeployedCode(data acceptance.TestData, tagSet string) string {
	tags := "{}"
	switch tagSet {
	case "single":
		tags = `{
    Environment = "Test"
  }`
	case "multiple":
		tags = `{
    Environment = "Updated"
    Label       = "EdgeAction"
  }`
	}

	return fmt.Sprintf(`
%[1]s

resource "azurerm_cdn_frontdoor_edge_action_version" "test" {
  name                         = "acctesteav%[2]d"
  cdn_frontdoor_edge_action_id = azurerm_cdn_frontdoor_edge_action.test.id
  deployment_type              = "file"
  code_file_name               = "edge_action.js"
  code_content_base64          = "%[3]s"

  tags = %[4]s
}

resource "azurerm_cdn_frontdoor_edge_action_version" "second" {
  name                         = "acctesteav2%[2]d"
  cdn_frontdoor_edge_action_id = azurerm_cdn_frontdoor_edge_action.test.id
  deployment_type              = "file"
  code_file_name               = "edge_action.js"
  code_content_base64          = "%[3]s"

  tags = %[4]s

  depends_on = [azurerm_cdn_frontdoor_edge_action_version.test]
}
`, r.template(data), data.RandomInteger, r.code("tags"), tags)
}
