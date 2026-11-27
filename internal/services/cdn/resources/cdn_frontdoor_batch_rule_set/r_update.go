// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_batch_rule_set

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/rulesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	cdnClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/client"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 4 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := cdnClient.From(metadata.Client).FrontDoorRuleSetsClient

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
