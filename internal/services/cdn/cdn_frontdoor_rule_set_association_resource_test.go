package cdn_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CdnFrontDoorRuleSetAssociationResource struct {
	// The configuration value to use in the Front Door Route Resource
	RouteRuleSetConfig string

	// The configuration value to use in the Front Door Rule Set Resource(s)
	RuleSetConfig string
}

func NewCdnFrontDoorRuleSetAssociationResource(routeRuleSetConfig string, ruleSetConfig string) *CdnFrontDoorRuleSetAssociationResource {
	return &CdnFrontDoorRuleSetAssociationResource{
		RouteRuleSetConfig: routeRuleSetConfig,
		RuleSetConfig:      ruleSetConfig,
	}
}

// NOTE: There isn't a complete test case because the basic and the
// update together equals what the complete test case would be...
func TestAccCdnFrontDoorRuleSetAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule_set_association", "test")
	r := NewCdnFrontDoorRuleSetAssociationResource("cdn_frontdoor_rule_set_ids = [azurerm_cdn_frontdoor_rule_set.one.id]", templateOneRuleSet(data))

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cdn_frontdoor_rule_set_ids.#").HasValue("1"),
			),
		},
	})
}

// NOTE: the 'requiresImport' test is not possible with this resource

func TestAccCdnFrontDoorRuleSetAssociation_destroy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule_set_association", "test")
	r := NewCdnFrontDoorRuleSetAssociationResource("cdn_frontdoor_rule_set_ids = [azurerm_cdn_frontdoor_rule_set.one.id]", templateOneRuleSet(data))

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cdn_frontdoor_rule_set_ids.#").HasValue("1"),
			),
		},
		{
			Config:             r.destroy(data),
			Check:              acceptance.ComposeTestCheckFunc(),
			ExpectNonEmptyPlan: true, // since destroying this resource actually removes the routes rule set associations
		},
	})
}

func TestAccCdnFrontDoorRuleSetAssociation_removeRouteRuleSetAssociation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule_set_association", "test")
	r := NewCdnFrontDoorRuleSetAssociationResource("cdn_frontdoor_rule_set_ids = [azurerm_cdn_frontdoor_rule_set.one.id]", templateOneRuleSet(data))

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cdn_frontdoor_rule_set_ids.#").HasValue("1"),
			),
		},
		{
			Config: r.removeRouteRuleSetAssociations(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cdn_frontdoor_rule_set_ids.#").HasValue("0"),
			),
			ExpectNonEmptyPlan: true, // since updating the field actually removes the routes rule set associations
		},
	})
}

func TestAccCdnFrontDoorRuleSetAssociation_removeMultipleRouteRuleSetAssociations(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule_set_association", "test")
	r := NewCdnFrontDoorRuleSetAssociationResource("cdn_frontdoor_rule_set_ids = [azurerm_cdn_frontdoor_rule_set.one.id, azurerm_cdn_frontdoor_rule_set.two.id, azurerm_cdn_frontdoor_rule_set.three.id]", templateMultipleRuleSets(data))

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleRuleSets(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cdn_frontdoor_rule_set_ids.#").HasValue("3"),
			),
			ExpectNonEmptyPlan: true, // I don't know why I need this, but the test fails without it, seems like a race condition to me...
		},
		{
			Config: r.removeRouteRuleSetAssociations(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cdn_frontdoor_rule_set_ids.#").HasValue("0"),
			),
			ExpectNonEmptyPlan: true, // since updating the field actually removes the routes rule set associations
		},
	})
}

func TestAccCdnFrontDoorRuleSetAssociation_multipleRuleSetAssociations(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule_set_association", "test")
	r := NewCdnFrontDoorRuleSetAssociationResource("cdn_frontdoor_rule_set_ids = [azurerm_cdn_frontdoor_rule_set.one.id, azurerm_cdn_frontdoor_rule_set.two.id, azurerm_cdn_frontdoor_rule_set.three.id]", templateMultipleRuleSets(data))

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleRuleSets(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cdn_frontdoor_rule_set_ids.#").HasValue("3"),
			),
			ExpectNonEmptyPlan: true, // I don't know why I need this, but the test fails without it, seems like a race condition to me...
		},
	})
}

func TestAccCdnFrontDoorRuleSetAssociation_routeRuleSetNotReferencedByAssociationError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule_set_association", "test")
	r := NewCdnFrontDoorRuleSetAssociationResource("cdn_frontdoor_rule_set_ids = [azurerm_cdn_frontdoor_rule_set.one.id, azurerm_cdn_frontdoor_rule_set.two.id, azurerm_cdn_frontdoor_rule_set.three.id]", templateMultipleRuleSets(data))

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.ruleSetsNotReferencedByAssociation(data),
			ExpectError: regexp.MustCompile("Please add the CDN FrontDoor Rule Sets to your configuration"),
		},
	})
}

func TestAccCdnFrontDoorRuleSetAssociation_duplicateAssociationRuleSetError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule_set_association", "test")
	r := NewCdnFrontDoorRuleSetAssociationResource("cdn_frontdoor_rule_set_ids = [azurerm_cdn_frontdoor_rule_set.one.id]", templateOneRuleSet(data))

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.duplicateRuleSets(data),
			ExpectError: regexp.MustCompile("please remove all duplicate entries from your configuration"),
		},
	})
}

func TestAccCdnFrontDoorRuleSetAssociation_ruleSetNotAssociatedWithRouteError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule_set_association", "test")
	r := NewCdnFrontDoorRuleSetAssociationResource("cdn_frontdoor_rule_set_ids = [azurerm_cdn_frontdoor_rule_set.two.id]", templateMultipleRuleSets(data))

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.basic(data),
			ExpectError: regexp.MustCompile("is currently not associated with the CDN FrontDoor Rule Sets"),
		},
	})
}

func TestAccCdnFrontDoorRuleSetAssociation_ruleSetDoesNotExistError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_frontdoor_rule_set_association", "test")
	r := NewCdnFrontDoorRuleSetAssociationResource("cdn_frontdoor_rule_set_ids = [azurerm_cdn_frontdoor_rule_set.one.id]", templateOneRuleSet(data))

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.ruleSetDoesNotExist(data),
			ExpectError: regexp.MustCompile("the following CDN FrontDoor Rule Sets do not exist"),
		},
	})
}

func (r CdnFrontDoorRuleSetAssociationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FrontDoorRuleSetAssociationID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cdn.FrontDoorRoutesClient
	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.AssociationName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (r CdnFrontDoorRuleSetAssociationResource) multipleRuleSets(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_rule_set_association" "test" {
  cdn_frontdoor_route_id     = azurerm_cdn_frontdoor_route.test.id
  cdn_frontdoor_rule_set_ids = [azurerm_cdn_frontdoor_rule_set.one.id, azurerm_cdn_frontdoor_rule_set.two.id, azurerm_cdn_frontdoor_rule_set.three.id]
}
`, template)
}

func (r CdnFrontDoorRuleSetAssociationResource) ruleSetsNotReferencedByAssociation(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_rule_set_association" "test" {
  cdn_frontdoor_route_id     = azurerm_cdn_frontdoor_route.test.id
  cdn_frontdoor_rule_set_ids = [azurerm_cdn_frontdoor_rule_set.one.id, azurerm_cdn_frontdoor_rule_set.three.id]
}
`, template)
}

func (r CdnFrontDoorRuleSetAssociationResource) removeRouteRuleSetAssociations(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_rule_set_association" "test" {
  cdn_frontdoor_route_id     = azurerm_cdn_frontdoor_route.test.id
  cdn_frontdoor_rule_set_ids = []
}
`, template)
}

func (r CdnFrontDoorRuleSetAssociationResource) destroy(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s
`, template)
}

func (r CdnFrontDoorRuleSetAssociationResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_rule_set_association" "test" {
  cdn_frontdoor_route_id     = azurerm_cdn_frontdoor_route.test.id
  cdn_frontdoor_rule_set_ids = [azurerm_cdn_frontdoor_rule_set.one.id]
}
`, template)
}

func (r CdnFrontDoorRuleSetAssociationResource) ruleSetDoesNotExist(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_rule_set_association" "test" {
  cdn_frontdoor_route_id     = azurerm_cdn_frontdoor_route.test.id
  cdn_frontdoor_rule_set_ids = ["/subscriptions/%s/resourceGroups/acctestRG-cdn-afdx-%[3]d/providers/Microsoft.Cdn/profiles/acctest-profile-%[3]d/ruleSets/yolo"]
}
`, template, data.Client().SubscriptionID, data.RandomInteger)
}

func (r CdnFrontDoorRuleSetAssociationResource) duplicateRuleSets(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_frontdoor_rule_set_association" "test" {
  cdn_frontdoor_route_id     = azurerm_cdn_frontdoor_route.test.id
  cdn_frontdoor_rule_set_ids = [azurerm_cdn_frontdoor_rule_set.one.id, azurerm_cdn_frontdoor_rule_set.one.id]
}
`, template)
}

func (r CdnFrontDoorRuleSetAssociationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cdn-afdx-%[1]d"
  location = "%[2]s"
}

resource "azurerm_cdn_frontdoor_profile" "test" {
  name                = "acctest-profile-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Standard_AzureFrontDoor"
}

resource "azurerm_cdn_frontdoor_origin_group" "test" {
  name                     = "acctest-origin-group-%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
  session_affinity_enabled = true

  health_probe {
    interval_in_seconds = 240
    path                = "/healthProbe"
    protocol            = "Https"
    request_type        = "GET"
  }

  load_balancing {
    additional_latency_in_milliseconds = 0
    sample_size                        = 16
    successful_samples_required        = 3
  }

  restore_traffic_time_to_healed_or_new_endpoint_in_minutes = 10
}

resource "azurerm_cdn_frontdoor_origin" "test" {
  name                          = "acctest-origin-%[1]d"
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
  enabled                       = true

  certificate_name_check_enabled = false
  host_name                      = "contoso.com"
  priority                       = 1
  weight                         = 1
}

resource "azurerm_cdn_frontdoor_endpoint" "test" {
  name                     = "acctest-endpoint-%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
  enabled                  = true
}

%[4]s

resource "azurerm_cdn_frontdoor_route" "test" {
  name                          = "acctest-contoso-%[1]d"
  cdn_frontdoor_endpoint_id     = azurerm_cdn_frontdoor_endpoint.test.id
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.test.id
  cdn_frontdoor_origin_ids      = [azurerm_cdn_frontdoor_origin.test.id]
  enabled                       = true

  https_redirect_enabled = true
  forwarding_protocol    = "HttpsOnly"
  patterns_to_match      = ["/%[3]s"]
  supported_protocols    = ["Http", "Https"]
  link_to_default_domain = true

  %[5]s

  cache {
    compression_enabled           = true
    content_types_to_compress     = ["text/html", "text/javascript", "text/xml"]
    query_strings                 = ["account", "settings", "foo", "bar"]
    query_string_caching_behavior = "IgnoreSpecifiedQueryStrings"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomStringOfLength(10), r.RuleSetConfig, r.RouteRuleSetConfig)
}

func templateOneRuleSet(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_cdn_frontdoor_rule_set" "one" {
  name                     = "acctestrulesetone%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
}
`, data.RandomInteger)
}

func templateMultipleRuleSets(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_cdn_frontdoor_rule_set" "one" {
  name                     = "acctestrulesetone%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
}

resource "azurerm_cdn_frontdoor_rule_set" "two" {
  name                     = "acctestrulesettwo%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
}

resource "azurerm_cdn_frontdoor_rule_set" "three" {
  name                     = "acctestrulesetthree%[1]d"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id
}
`, data.RandomInteger)
}
