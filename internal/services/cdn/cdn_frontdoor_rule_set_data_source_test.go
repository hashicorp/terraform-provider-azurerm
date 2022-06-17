package cdn_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type CdnFrontDoorRuleSetDataSource struct{}

func TestAccCdnFrontDoorRuleSetDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cdn_frontdoor_rule_set", "test")
	d := CdnFrontDoorRuleSetDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("cdn_frontdoor_profile_id").MatchesOtherKey(check.That("azurerm_cdn_frontdoor_profile.test").Key("id")),
			),
		},
	})
}

func (CdnFrontDoorRuleSetDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_cdn_frontdoor_rule_set" "test" {
  name                = azurerm_cdn_frontdoor_rule_set.test.name
  profile_name        = azurerm_cdn_frontdoor_rule_set.test.profile_name
  resource_group_name = azurerm_cdn_frontdoor_rule_set.test.resource_group_name
}
`, CdnFrontDoorRuleSetResource{}.complete(data))
}
