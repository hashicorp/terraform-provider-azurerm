// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_batch_rule_set

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rulesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 4 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cdn.FrontDoorRuleSetsClient

			id, err := rulesets.ParseRuleSetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model CdnFrontDoorBatchRuleSetModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			payload, err := expandCdnFrontDoorBatchRuleSetPayload(model)
			if err != nil {
				return err
			}

			if err := client.CreateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
