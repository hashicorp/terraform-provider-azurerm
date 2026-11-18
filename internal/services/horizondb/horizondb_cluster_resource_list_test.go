// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package horizondb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/querycheck/queryfilter"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
	"github.com/zclconf/go-cty/cty"
)

func TestAccHorizondbCluster_list(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_horizondb_cluster", "test")
	r := HorizondbClusterResource{}
	checks := make([]querycheck.QueryResultCheck, 0, 2)
	resourceChecks := make([]querycheck.QueryResultCheck, 0, 2)
	for index := range 2 {
		identity := map[string]knownvalue.Check{
			"name":                knownvalue.StringExact(fmt.Sprintf("acctest-horizondb-%d-%d", data.RandomInteger, index)),
			"resource_group_name": knownvalue.StringExact(fmt.Sprintf("acctestRG-horizondb-%d", data.RandomInteger)),
			"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
		}
		checks = append(checks, querycheck.ExpectIdentity("azurerm_horizondb_cluster.list", identity))
		resourceChecks = append(resourceChecks, querycheck.ExpectResourceKnownValues("azurerm_horizondb_cluster.list",
			queryfilter.ByResourceIdentity(identity), []querycheck.KnownValueCheck{{
				Path:       tfjsonpath.New("postgresql_version"),
				KnownValue: knownvalue.StringExact("18"),
			}, {
				Path:       tfjsonpath.New("read_write_pool").AtSliceIndex(0).AtMapKey("name"),
				KnownValue: knownvalue.StringExact("primary"),
			}, {
				Path:       tfjsonpath.New("read_write_pool").AtSliceIndex(0).AtMapKey("id"),
				KnownValue: knownvalue.StringExact(fmt.Sprintf("/subscriptions/%s/resourceGroups/acctestRG-horizondb-%d/providers/Microsoft.HorizonDb/clusters/acctest-horizondb-%d-%d/pools/primary", data.Subscriptions.Primary, data.RandomInteger, data.RandomInteger, index)),
			}}))
	}
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.SkipBelow(tfversion.Version1_14_0)},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{Config: r.listConfig(data)},
			{
				Query: true,
				Config: `
list "azurerm_horizondb_cluster" "list" {
  provider = azurerm
  config {}
}
`,
				QueryResultChecks: append(checks, querycheck.ExpectLengthAtLeast("azurerm_horizondb_cluster.list", 2)),
			},
			{
				Query: true,
				Config: `
list "azurerm_horizondb_cluster" "list" {
  provider         = azurerm
  include_resource = true
  config {
    resource_group_name = azurerm_resource_group.test.name
  }
}
`,
				QueryResultChecks: append(append(resourceChecks, checks...), querycheck.ExpectLength("azurerm_horizondb_cluster.list", 2)),
			},
		},
	})
}

func TestHorizondbClusterListConfig(t *testing.T) {
	data := acceptance.TestData{RandomInteger: 12345}
	data.Locations.Primary = "eastus2"
	file, diags := hclsyntax.ParseConfig([]byte(HorizondbClusterResource{}.listConfig(data)), "list.tf", hcl.InitialPos)
	if diags.HasErrors() {
		t.Fatalf("parsing list setup: %s", diags.Error())
	}
	for _, block := range file.Body.(*hclsyntax.Body).Blocks {
		if block.Type != "resource" || len(block.Labels) != 2 || block.Labels[0] != "azurerm_horizondb_cluster" || block.Labels[1] != "test" {
			continue
		}
		for name, expected := range map[string]string{
			"administrator_login":    "acctestadmin",
			"administrator_password": "P@ssword-123!",
		} {
			attribute, ok := block.Body.Attributes[name]
			if !ok {
				t.Errorf("list setup is missing required attribute %q", name)
				continue
			}
			value, diags := attribute.Expr.Value(nil)
			if diags.HasErrors() {
				t.Errorf("evaluating %q: %s", name, diags.Error())
				continue
			}
			if !value.RawEquals(cty.StringVal(expected)) {
				t.Errorf("list setup has an unexpected value for %q", name)
			}
		}
		return
	}
	t.Fatal("list setup is missing the cluster resource")
}

func (r HorizondbClusterResource) listConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_horizondb_cluster" "test" {
  count                  = 2
  name                   = "acctest-horizondb-%d-${count.index}"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  postgresql_version     = "18"
  administrator_login    = "acctestadmin"
  administrator_password = "P@ssword-123!"

  read_write_pool {
    name                   = "primary"
    node_count             = 1
    high_availability_mode = "Disabled"

    compute {
      processor  = "AMD"
      generation = 1
      vcores     = 4
    }
  }
}
`, r.template(data), data.RandomInteger)
}
