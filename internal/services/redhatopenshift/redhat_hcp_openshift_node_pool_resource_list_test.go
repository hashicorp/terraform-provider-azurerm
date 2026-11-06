// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redhatopenshift_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/querycheck/queryfilter"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccRedhatHcpOpenshiftNodePool_list_basic(t *testing.T) {
	r := RedhatHcpOpenshiftNodePoolResource{}
	data := acceptance.BuildTestData(t, "azurerm_redhat_hcp_openshift_node_pool", "test")
	listResourceAddress := "azurerm_redhat_hcp_openshift_node_pool.list"

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basic(data),
			},
			{
				Query:  true,
				Config: r.basicListQuery(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast(listResourceAddress, 1),
					querycheck.ExpectIdentity(listResourceAddress, map[string]knownvalue.Check{
						"name":                        knownvalue.StringExact("acctestnp" + data.RandomString),
						"hcp_open_shift_cluster_name": knownvalue.StringExact(fmt.Sprintf("acctestarohcp%d", data.RandomInteger)),
						"resource_group_name":         knownvalue.StringExact(fmt.Sprintf("acctestRG-aro-hcp-%d", data.RandomInteger)),
						"subscription_id":             knownvalue.StringExact(data.Subscriptions.Primary),
					}),
					querycheck.ExpectResourceKnownValues(listResourceAddress, queryfilter.ByDisplayName(knownvalue.StringExact("acctestnp"+data.RandomString)), []querycheck.KnownValueCheck{
						{
							Path:       tfjsonpath.New("hcp_openshift_cluster_id"),
							KnownValue: knownvalue.StringExact(fmt.Sprintf("/subscriptions/%s/resourceGroups/acctestRG-aro-hcp-%d/providers/Microsoft.RedHatOpenShift/hcpOpenShiftClusters/acctestarohcp%d", data.Subscriptions.Primary, data.RandomInteger, data.RandomInteger)),
						},
						{
							Path:       tfjsonpath.New("replicas"),
							KnownValue: knownvalue.Int64Exact(1),
						},
						{
							Path:       tfjsonpath.New("platform").AtSliceIndex(0).AtMapKey("vm_size"),
							KnownValue: knownvalue.StringExact("Standard_D4s_v3"),
						},
					}),
				},
			},
			{
				Config: r.removeFromState(RedhatHcpOpenshiftClusterResource{}.basic(data)),
			},
		},
	})
}

func (r RedhatHcpOpenshiftNodePoolResource) basicListQuery(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_redhat_hcp_openshift_node_pool" "list" {
  provider         = azurerm
  include_resource = true
  config {
    hcp_openshift_cluster_id = "/subscriptions/%[1]s/resourceGroups/acctestRG-aro-hcp-%[2]d/providers/Microsoft.RedHatOpenShift/hcpOpenShiftClusters/acctestarohcp%[2]d"
  }
}
`, data.Subscriptions.Primary, data.RandomInteger)
}
