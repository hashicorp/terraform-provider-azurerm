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

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.PrefixListResources
			rulestackClient := metadata.Client.PaloAlto.LocalRulestackResources
			model := LocalRuleStackPrefixListModel{}

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			rulestackId, err := localrulestackresources.ParseLocalRulestackID(model.RuleStackID)
			if err != nil {
				return err
			}
			locks.ByID(rulestackId.ID())
			defer locks.UnlockByID(rulestackId.ID())

			id := prefixlistresources.NewLocalRulestackPrefixListID(rulestackId.SubscriptionId, rulestackId.ResourceGroupName, rulestackId.LocalRulestackName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.PrefixListLocalRulestackGet(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			props := prefixlistresources.PrefixObject{
				PrefixList: model.PrefixList,
			}

			if model.AuditComment != "" {
				props.AuditComment = pointer.To(model.AuditComment)
			}

			if model.Description != "" {
				props.Description = pointer.To(model.Description)
			}

			prefixList := prefixlistresources.PrefixListResource{
				Properties: props,
			}

			if err = client.PrefixListLocalRulestackCreateOrUpdateCallbackThenPoll(ctx, id, prefixList, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			if err = rulestackClient.LocalRulestackscommitThenPoll(ctx, *rulestackId); err != nil {
				return fmt.Errorf("committing Local RuleStack config for %s: %+v", id, err)
			}

			return nil
		},
	}
}
