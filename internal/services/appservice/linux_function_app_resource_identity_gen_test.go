package appservice_test

import (
	"context"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccLinuxFunctionApp_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_function_app", "test")
	r := LinuxFunctionAppResource{}

	resource.ParallelTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.12.0-rc2"))),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basic(data, "B1"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectIdentity("azurerm_linux_function_app.test", map[string]knownvalue.Check{
						"kind":            knownvalue.StringExact("functionapp,linux"),
						"subscription_id": knownvalue.StringExact(data.Subscriptions.Primary), // The identity has 4 properties, but we can't check them all here. statechecks for paths are not knownvalue.Checks?
					}),
					statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_linux_function_app.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("resource_group_name")),
					statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_linux_function_app.test", tfjsonpath.New("site_name"), tfjsonpath.New("name")),
				},
			},
			data.ImportBlockWithResourceIdentityStep(),
			data.ImportBlockWithIDStep(),
		},
	})
}
