// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"context"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccDedicatedHost_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dedicated_host", "test")
	r := DedicatedHostResource{}

	resource.ParallelTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.12.0-rc2"))),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basic(data),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_dedicated_host.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
					customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_dedicated_host.test", tfjsonpath.New("host_group_name"), tfjsonpath.New("dedicated_host_group_id")),
					customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_dedicated_host.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("dedicated_host_group_id")),
					customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_dedicated_host.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("dedicated_host_group_id")),
				},
			},
			data.ImportBlockWithResourceIdentityStep(),
			data.ImportBlockWithIDStep(),
		},
	})
}
