// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

var cases = map[string][]string{
	"test-1": {"/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/resourcegroups/resGroup1/providers/Microsoft.web/siTes/site1", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1"},
	"test-2": {"/subscriptions/12345678-1234-9876-4563-123456789012/ResourceGroups/resGroup1/PROVIDERS/Microsoft.APIManagement/service/service1/gateWays/gateway1/hostnameconfigurations/config1", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ApiManagement/service/service1/gateways/gateway1/hostnameConfigurations/config1"},
	"test-3": {"/SubScripTionS/12345678-1234-9876-4563-123456789012/resourceGROUPS/resGroup1/providers/microsoft.apiManagement/Service/service1/apis/api1", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ApiManagement/service/service1/apis/api1"},
	"test-4": {"/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualmachinescalesets/scaleSet1/virtualmachines/machine1/networkinterfaCes/networkInterface1/ipconFigurations/ipConfig1/PublicipAddresses/publicAddress1", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/scaleSet1/virtualMachines/machine1/networkInterfaces/networkInterface1/ipConfigurations/ipConfig1/publicIPAddresses/publicAddress1"},
	"test-5": {"/subscriptions/12345678-1234-9876-4563-123456789012/resourcegroups/resGroup1/providers/microsoft.chaos/Targets/target1", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Chaos/targets/target1"},
	"test-6": {"SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/resourcegroups/resGroup1/providers/Microsoft.web/siTes/site1", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1"},
}

func TestProviderFunctionNormaliseResourceID_multiple(t *testing.T) {
	if !features.FourPointOhBeta() {
		t.Skipf("skipping test due to missing feature flag")
	}
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.8.0-beta1"))),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: testOutputMultiple(cases),
				Check: acceptance.ComposeTestCheckFunc(
					resource.TestCheckOutput("test-1", cases["test-1"][1]),
					resource.TestCheckOutput("test-2", cases["test-2"][1]),
					resource.TestCheckOutput("test-3", cases["test-3"][1]),
					resource.TestCheckOutput("test-4", cases["test-4"][1]),
					resource.TestCheckOutput("test-5", cases["test-5"][1]),
				),
			},
		},
	})
}

func testOutputMultiple(cases map[string][]string) string {
	outputs := ""
	for k, v := range cases {
		outputs += fmt.Sprintf(`

output "%s" {
  value = provider::azurerm::normalise_resource_id("%s")
}

`, k, v[0])
	}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s
`, outputs)
}
