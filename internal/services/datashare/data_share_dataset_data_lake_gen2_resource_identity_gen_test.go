// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package datashare_test

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
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccDataShareDatasetDataLakeGen2_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_dataset_data_lake_gen_2", "test")
	r := DataShareDatasetDataLakeGen2Resource{}

	resource.ParallelTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.12.0-rc2"))),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicFile(data),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectIdentityValue("azurerm_data_share_dataset_data_lake_gen_2.test", tfjsonpath.New("subscription_id"), knownvalue.StringExact(data.Subscriptions.Primary)),
					statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_data_share_dataset_data_lake_gen_2.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
					customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_data_share_dataset_data_lake_gen_2.test", tfjsonpath.New("account_name"), tfjsonpath.New("share_id")),
					customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_data_share_dataset_data_lake_gen_2.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("share_id")),
					customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_data_share_dataset_data_lake_gen_2.test", tfjsonpath.New("share_name"), tfjsonpath.New("share_id")),
				},
			},
			data.ImportBlockWithResourceIdentityStep(),
			data.ImportBlockWithIDStep(),
		},
	})
}
