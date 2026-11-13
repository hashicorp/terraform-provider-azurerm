// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_metadata_schema_test

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

func TestAccApiCenterMetadataSchema_listByServiceID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_center_metadata_schema", "list")
	r := ApiCenterMetadataSchemaResource{}

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
					querycheck.ExpectLengthAtLeast("azurerm_api_center_metadata_schema.list", 1),
					querycheck.ExpectIdentity(
						"azurerm_api_center_metadata_schema.list",
						map[string]knownvalue.Check{
							"name":                knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"service_name":        knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
		},
	})
}

func (r ApiCenterMetadataSchemaResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_center_metadata_schema" "test" {
  count      = 3
  name       = "acctest%d${count.index}"
  service_id = azurerm_api_center_service.test.id

  metadata_assignment {
    entity = "api"
  }

  schema = jsonencode({
    type = "object"
    properties = {
      name = {
        type = "string"
      }
    }
  })
}
`, r.template(data), data.RandomInteger)
}

func (r ApiCenterMetadataSchemaResource) basicQuery() string {
	return `
list "azurerm_api_center_metadata_schema" "list" {
  provider = azurerm
  config {
    api_center_service_id = azurerm_api_center_service.test.id
  }
}
`
}
