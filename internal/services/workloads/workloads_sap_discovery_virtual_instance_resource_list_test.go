// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package workloads_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func testAccWorkloadsSAPDiscoveryVirtualInstance_list_basic(t *testing.T) {
	// The dependent central server VM requires many complicated manual configurations. So it has to test based on the resource provided by service team.
	if os.Getenv("ARM_TEST_SAP_VIRTUAL_INSTANCE_NAME") == "" || os.Getenv("ARM_TEST_CENTRAL_SERVER_VM_ID") == "" || os.Getenv("ARM_TEST_IDENTITY_ID") == "" {
		t.Skip("Skipping as `ARM_TEST_SAP_VIRTUAL_INSTANCE_NAME`, `ARM_TEST_CENTRAL_SERVER_VM_ID` and `ARM_TEST_IDENTITY_ID` are not specified")
	}

	r := WorkloadsSapDiscoveryVirtualInstanceResource{}
	listResourceAddress := "azurerm_workloads_sap_discovery_virtual_instance.list"
	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_discovery_virtual_instance", "test")

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

func (r WorkloadsSapDiscoveryVirtualInstanceResource) basicListQuery() string {
	return `
list "azurerm_workloads_sap_discovery_virtual_instance" "list" {
  provider = azurerm
  config {}
}
`
}

func (r WorkloadsSapDiscoveryVirtualInstanceResource) basicListQueryByResourceGroupName(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_workloads_sap_discovery_virtual_instance" "list" {
  provider = azurerm
  config {
    resource_group_name = "acctestRG-sapvis-%[1]d"
  }
}
`, data.RandomInteger)
}
