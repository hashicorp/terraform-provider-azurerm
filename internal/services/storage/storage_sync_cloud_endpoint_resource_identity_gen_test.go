// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage_test

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

func TestAccStorageSyncCloudEndpoint_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_sync_cloud_endpoint", "test")
	r := StorageSyncCloudEndpointResource{}

	resource.ParallelTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.12.0-rc2"))),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basic(data),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_storage_sync_cloud_endpoint.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
					customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_storage_sync_cloud_endpoint.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("storage_sync_group_id")),
					customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_storage_sync_cloud_endpoint.test", tfjsonpath.New("storage_sync_service_name"), tfjsonpath.New("storage_sync_group_id")),
					customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_storage_sync_cloud_endpoint.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("storage_sync_group_id")),
					customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_storage_sync_cloud_endpoint.test", tfjsonpath.New("sync_group_name"), tfjsonpath.New("storage_sync_group_id")),
				},
			},
			data.ImportBlockWithResourceIdentityStep(false),
			data.ImportBlockWithIDStep(false),
		},
	})
}
