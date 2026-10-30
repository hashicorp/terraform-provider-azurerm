// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_local_rulestack_fqdn_list

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/fqdnlistlocalrulestackresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/localrulestackresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.FqdnListLocalRulestackResources
			rulestackClient := metadata.Client.PaloAlto.LocalRulestackResources

			model := LocalRulestackFQDNListModel{}

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id, err := fqdnlistlocalrulestackresources.ParseLocalRulestackFqdnListID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			rulestackId := localrulestackresources.NewLocalRulestackID(id.SubscriptionId, id.ResourceGroupName, id.LocalRulestackName)
			locks.ByID(rulestackId.ID())
			defer locks.UnlockByID(rulestackId.ID())

			existing, err := client.FqdnListLocalRulestackGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			fqdnList := *existing.Model

			if metadata.ResourceData.HasChange("fully_qualified_domain_names") {
				fqdnList.Properties.FqdnList = model.FQDNList
			}

			if metadata.ResourceData.HasChange("audit_comment") {
				fqdnList.Properties.AuditComment = pointer.To(model.AuditComment)
			}

			if metadata.ResourceData.HasChange("description") {
				fqdnList.Properties.Description = pointer.To(model.Description)
			}

			if _, err = client.FqdnListLocalRulestackCreateOrUpdate(ctx, *id, fqdnList); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			if err = rulestackClient.LocalRulestackscommitThenPoll(ctx, rulestackId); err != nil {
				return fmt.Errorf("committing Local Rulestack config for %s: %+v", id, err)
			}

			return nil
		},
	}
}
