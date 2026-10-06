// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_test

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
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage"
)

func TestAccStorageContainer_list(t *testing.T) {
	data := acceptance.BuildTestData(t, storage.StorageContainerResourceName, "list")
	r := StorageContainerResource{}

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
					querycheck.ExpectLength(data.ResourceName, 3),
					querycheck.ExpectIdentity(
						data.ResourceName,
						map[string]knownvalue.Check{
							"name":                 knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"storage_account_name": knownvalue.StringRegexp(regexp.MustCompile(data.RandomString)),
							"resource_group_name":  knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":      knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
		},
	})
}

func (r StorageContainerResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_container" "test" {
  count = 3

  name               = "acctest-container-${count.index}-%[2]d"
  storage_account_id = azurerm_storage_account.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r StorageContainerResource) basicQuery() string {
	return `
list "azurerm_storage_container" "list" {
  provider = azurerm
  config {
    storage_account_id = azurerm_storage_container.test[0].storage_account_id
  }
}
`
}
