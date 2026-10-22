// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_local_rulestack_prefix_list

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/prefixlistresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.PrefixListResources

			id, err := prefixlistresources.ParseLocalRulestackPrefixListID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state LocalRuleStackPrefixListModel

			existing, err := client.PrefixListLocalRulestackGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			state.Name = id.PrefixListName
			state.RuleStackID = prefixlistresources.NewLocalRulestackID(id.SubscriptionId, id.ResourceGroupName, id.LocalRulestackName).ID()
			if model := existing.Model; model != nil {
				props := model.Properties

				state.PrefixList = props.PrefixList
				state.AuditComment = pointer.From(props.AuditComment)
				state.Description = pointer.From(props.Description)
			}

			return metadata.Encode(&state)
		},
	}
}
