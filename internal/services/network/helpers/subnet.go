// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/subnets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func SubnetProvisioningStateRefreshFunc(ctx context.Context, client *subnets.SubnetsClient, id commonids.SubnetId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		res, err := client.Get(ctx, id, subnets.DefaultGetOperationOptions())
		if err != nil {
			return nil, "", fmt.Errorf("polling for %s: %+v", id.String(), err)
		}

		if res.Model != nil && res.Model.Properties != nil && res.Model.Properties.ProvisioningState != nil {
			return res, string(*res.Model.Properties.ProvisioningState), nil
		}
		return nil, "", fmt.Errorf("unable to read provisioning state")
	}
}
