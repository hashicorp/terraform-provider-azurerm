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
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.RaiPoliciesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model AccountRaiPolicyResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			cognitiveAccountId, err := raipolicies.ParseAccountID(model.AccountId)
			if err != nil {
				return err
			}

			id := raipolicies.NewAccountRaiPolicyID(subscriptionId, cognitiveAccountId.ResourceGroupName, cognitiveAccountId.AccountName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			locks.ByID(cognitiveAccountId.ID())
			defer locks.UnlockByID(cognitiveAccountId.ID())

			raiPolicy := raipolicies.RaiPolicy{
				Name: pointer.To(model.Name),
				Properties: &raipolicies.RaiPolicyProperties{
					BasePolicyName: pointer.To(model.BasePolicyName),
					ContentFilters: expandRaiPolicyContentFilters(model.ContentFilter),
				},
				Tags: pointer.To(model.Tags),
			}

			if model.Mode != "" {
				raiPolicy.Properties.Mode = pointer.ToEnum[raipolicies.RaiPolicyMode](model.Mode)
			}

			if _, err := client.CreateOrUpdate(ctx, id, raiPolicy); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func expandRaiPolicyContentFilters(filters []AccountRaiPolicyContentFilter) *[]raipolicies.RaiPolicyContentFilter {
	if filters == nil {
		return nil
	}

	contentFilters := make([]raipolicies.RaiPolicyContentFilter, 0, len(filters))
	for _, filter := range filters {
		f := raipolicies.RaiPolicyContentFilter{
			Name:     pointer.To(filter.Name),
			Enabled:  pointer.To(filter.FilterEnabled),
			Blocking: pointer.To(filter.BlockEnabled),
			Source:   pointer.ToEnum[raipolicies.RaiPolicyContentSource](filter.Source),
		}

		if filter.SeverityThreshold != "" {
			f.SeverityThreshold = pointer.ToEnum[raipolicies.ContentLevel](filter.SeverityThreshold)
		}

		contentFilters = append(contentFilters, f)
	}
	return &contentFilters
}
