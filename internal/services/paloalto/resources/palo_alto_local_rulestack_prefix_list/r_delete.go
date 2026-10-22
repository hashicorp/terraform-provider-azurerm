// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_local_rulestack_prefix_list

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/localrulestackresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/prefixlistresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.PrefixListResources
			rulestackClient := metadata.Client.PaloAlto.LocalRulestackResources

			id, err := prefixlistresources.ParseLocalRulestackPrefixListID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			rulestackId := localrulestackresources.NewLocalRulestackID(id.SubscriptionId, id.ResourceGroupName, id.LocalRulestackName)
			locks.ByID(rulestackId.ID())
			defer locks.UnlockByID(rulestackId.ID())

			if err = client.PrefixListLocalRulestackDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			if err = rulestackClient.LocalRulestackscommitThenPoll(ctx, rulestackId); err != nil {
				return fmt.Errorf("committing Local Rulestack config for %s: %+v", id, err)
			}

			return nil
		},
	}
}
