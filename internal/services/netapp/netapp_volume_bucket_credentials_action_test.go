// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

type NetAppVolumeBucketCredentialsAction struct{}

func TestAccNetAppVolumeBucketCredentialsAction_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_bucket_credentials", "test")
	a := NetAppVolumeBucketCredentialsAction{}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: a.basic(data),
				Check:  nil,
			},
		},
	})
}

func (a NetAppVolumeBucketCredentialsAction) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

action "azurerm_netapp_volume_bucket_credentials" "test" {
  config {
    bucket_id            = azurerm_netapp_volume_bucket.test.id
    key_pair_expiry_days = 30
  }
}

resource "terraform_data" "trigger" {
  input = azurerm_netapp_volume_bucket.test.id
  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.azurerm_netapp_volume_bucket_credentials.test]
    }
  }
}
`, NetAppVolumeBucketResource{}.withKeyVault(data))
}
