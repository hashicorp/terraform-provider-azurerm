// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package redhatopenshift_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccRedhatOpenshiftCluster_list_basic(t *testing.T) {
	r := RedhatOpenshiftClusterResource{}
	listResourceAddress := "azurerm_redhat_openshift_cluster.list"
	data := acceptance.BuildTestData(t, "azurerm_redhat_openshift_cluster", "test")

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
				Config: r.basicListQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast(listResourceAddress, 1),
				},
			},
			{
				Query:  true,
				Config: r.basicListQueryByResourceGroupName(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast(listResourceAddress, 1),
				},
			},
		},
	})
}

func (r RedhatOpenshiftClusterResource) basicListQuery() string {
	return `
list "azurerm_redhat_openshift_cluster" "list" {
  provider = azurerm
  config {}
}
`
}

func (r RedhatOpenshiftClusterResource) basicListQueryByResourceGroupName(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_redhat_openshift_cluster" "list" {
  provider = azurerm
  config {
    resource_group_name = "acctestRG-aro-%[1]d"
  }
}
`, data.RandomInteger)
}
