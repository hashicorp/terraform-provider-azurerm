// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_rule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (c Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 6 * time.Hour,
		Func: func(ctx context.Context, rmd sdk.ResourceMetaData) error {
			client := rmd.Client.Cdn.FrontDoorRulesClient

			id, err := rules.ParseRuleID(rmd.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
