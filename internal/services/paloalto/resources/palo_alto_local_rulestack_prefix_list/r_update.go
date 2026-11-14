// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_local_rulestack_prefix_list

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/localrulestackresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/prefixlistresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.PrefixListResources
			rulestackClient := metadata.Client.PaloAlto.LocalRulestackResources

			id, err := prefixlistresources.ParseLocalRulestackPrefixListID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			model := LocalRuleStackPrefixListModel{}

			if err = metadata.Decode(&model); err != nil {
				return err
			}

			rulestackId := localrulestackresources.NewLocalRulestackID(id.SubscriptionId, id.ResourceGroupName, id.LocalRulestackName)
			locks.ByID(rulestackId.ID())
			defer locks.UnlockByID(rulestackId.ID())

			existing, err := client.PrefixListLocalRulestackGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s for update: %+v", *id, err)
			}

			prefixList := *existing.Model

			if metadata.ResourceData.HasChange("prefix_list") {
				prefixList.Properties.PrefixList = model.PrefixList
			}

			if metadata.ResourceData.HasChange("audit_comment") {
				prefixList.Properties.AuditComment = pointer.To(model.AuditComment)
			}

			if metadata.ResourceData.HasChange("description") {
				prefixList.Properties.Description = pointer.To(model.Description)
			}

			if _, err = client.PrefixListLocalRulestackCreateOrUpdate(ctx, *id, prefixList); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			if err = rulestackClient.LocalRulestackscommitThenPoll(ctx, rulestackId); err != nil {
				return fmt.Errorf("committing Local Rulestack config for %s: %+v", id, err)
			}

			return nil
		},
	}
}
