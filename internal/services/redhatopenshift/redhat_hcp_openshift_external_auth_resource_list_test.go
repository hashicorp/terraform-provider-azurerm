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

func TestAccRedhatHcpOpenshiftExternalAuth_list_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_hcp_openshift_external_auth", "test")
	testResource := RedhatHcpOpenshiftExternalAuthResource{}
	listResourceAddress := "azurerm_redhat_hcp_openshift_external_auth.list"
	name := "acctestea" + data.RandomString

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: testResource.basic(data),
			},
			{
				Query:  true,
				Config: testResource.basicListQuery(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast(listResourceAddress, 1),
					querycheck.ExpectIdentity(listResourceAddress, map[string]knownvalue.Check{
						"name":                        knownvalue.StringExact(name),
						"hcp_open_shift_cluster_name": knownvalue.StringExact(fmt.Sprintf("acctestarohcp%d", data.RandomInteger)),
						"resource_group_name":         knownvalue.StringExact(fmt.Sprintf("acctestRG-aro-hcp-%d", data.RandomInteger)),
						"subscription_id":             knownvalue.StringExact(data.Subscriptions.Primary),
					}),
					querycheck.ExpectResourceKnownValues(listResourceAddress, queryfilter.ByDisplayName(knownvalue.StringExact(name)), []querycheck.KnownValueCheck{
						{
							Path:       tfjsonpath.New("hcp_openshift_cluster_id"),
							KnownValue: knownvalue.StringExact(fmt.Sprintf("/subscriptions/%s/resourceGroups/acctestRG-aro-hcp-%d/providers/Microsoft.RedHatOpenShift/hcpOpenShiftClusters/acctestarohcp%d", data.Subscriptions.Primary, data.RandomInteger, data.RandomInteger)),
						},
						{
							Path:       tfjsonpath.New("claim").AtSliceIndex(0).AtMapKey("mappings").AtSliceIndex(0).AtMapKey("username").AtSliceIndex(0).AtMapKey("claim"),
							KnownValue: knownvalue.StringExact("sub"),
						},
					}),
				},
			},
		},
	})
}

func (RedhatHcpOpenshiftExternalAuthResource) basicListQuery(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_redhat_hcp_openshift_external_auth" "list" {
  provider         = azurerm
  include_resource = true
  config {
    hcp_openshift_cluster_id = "/subscriptions/%[1]s/resourceGroups/acctestRG-aro-hcp-%[2]d/providers/Microsoft.RedHatOpenShift/hcpOpenShiftClusters/acctestarohcp%[2]d"
  }
}
`, data.Subscriptions.Primary, data.RandomInteger)
}
