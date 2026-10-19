// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_rai_policy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/raipolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.RaiPoliciesClient

			id, err := raipolicies.ParseAccountRaiPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			cognitiveAccountId := raipolicies.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName)

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := AccountRaiPolicyResourceModel{
				Name:      id.RaiPolicyName,
				AccountId: cognitiveAccountId.ID(),
			}

			if model := resp.Model; model != nil {
				state.Tags = pointer.From(model.Tags)

				if props := model.Properties; props != nil {
					state.BasePolicyName = pointer.From(props.BasePolicyName)
					state.ContentFilter = flattenRaiPolicyContentFilters(props.ContentFilters)
					state.Mode = pointer.FromEnum(props.Mode)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenRaiPolicyContentFilters(filters *[]raipolicies.RaiPolicyContentFilter) []AccountRaiPolicyContentFilter {
	contentFilters := make([]AccountRaiPolicyContentFilter, 0)
	if filters == nil {
		return contentFilters
	}

	for _, filter := range *filters {
		contentFilters = append(contentFilters, AccountRaiPolicyContentFilter{
			Name:              pointer.From(filter.Name),
			FilterEnabled:     pointer.From(filter.Enabled),
			BlockEnabled:      pointer.From(filter.Blocking),
			SeverityThreshold: pointer.FromEnum(filter.SeverityThreshold),
			Source:            pointer.FromEnum(filter.Source),
		})
	}
	return contentFilters
}

type AccountRaiPolicyContentFilter struct {
	Name              string `tfschema:"name"`
	FilterEnabled     bool   `tfschema:"filter_enabled"`
	BlockEnabled      bool   `tfschema:"block_enabled"`
	SeverityThreshold string `tfschema:"severity_threshold"`
	Source            string `tfschema:"source"`
}

var severityThresholdNotApplicableFilterNames = []string{
	"Jailbreak",
	"Indirect Attack",
	"Protected Material Text",
	"Protected Material Code",
}
