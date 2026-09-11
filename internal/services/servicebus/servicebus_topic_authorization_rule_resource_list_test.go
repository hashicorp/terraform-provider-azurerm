package servicebus_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccServiceBusTopicAuthorizationRule_listByTopicID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_servicebus_topic_authorization_rule", "testlist1")
	r := ServicebusTopicAuthorizationRuleResource{}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicList(data),
			},
			{
				Query:  true,
				Config: r.basicQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("azurerm_servicebus_topic_authorization_rule.list", 3),
					querycheck.ExpectIdentity(
						"azurerm_servicebus_topic_authorization_rule.list",
						map[string]knownvalue.Check{
							"name":                knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"topic_name":          knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
							"namespace_name":      knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
						},
					),
				},
			},
		},
	})
}

func (r ServicebusTopicAuthorizationRuleResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctest-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "test" {
  name         = "acctestservicebustopic-%[1]d"
  namespace_id = azurerm_servicebus_namespace.test.id
}

resource "azurerm_servicebus_topic_authorization_rule" "test" {
  count = 3

  name     = "acctest-%[1]d-${count.index}"
  topic_id = azurerm_servicebus_topic.test.id
  listen   = true
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ServicebusTopicAuthorizationRuleResource) basicQuery() string {
	return `
list "azurerm_servicebus_topic_authorization_rule" "list" {
  provider = azurerm
  config {
    servicebus_topic_id = azurerm_servicebus_topic.test.id
  }
}
`
}
