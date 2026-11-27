// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network_security_perimeter_access_rule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecurityperimeteraccessrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	networkmanagerClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/networkmanager/client"
)

func (Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := networkmanagerClient.From(metadata.Client).NetworkSecurityPerimeterAccessRulesClient

			id, err := networksecurityperimeteraccessrules.ParseAccessRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			return nil
		},
	}
}
