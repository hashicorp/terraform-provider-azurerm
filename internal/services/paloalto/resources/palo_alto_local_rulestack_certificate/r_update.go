// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_local_rulestack_certificate

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/certificateobjectlocalrulestackresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/localrulestackresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.CertificateObjectLocalRulestackResources
			rulestackClient := metadata.Client.PaloAlto.LocalRulestackResources
			model := LocalRuleStackCertificateModel{}

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id, err := certificateobjectlocalrulestackresources.ParseLocalRulestackCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			rulestackId := localrulestackresources.NewLocalRulestackID(id.SubscriptionId, id.ResourceGroupName, id.LocalRulestackName)
			locks.ByID(rulestackId.ID())
			defer locks.UnlockByID(rulestackId.ID())

			existing, err := client.CertificateObjectLocalRulestackGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			cert := *existing.Model

			if metadata.ResourceData.HasChange("description") {
				cert.Properties.Description = pointer.To(model.Description)
			}

			if metadata.ResourceData.HasChange("audit_comment") {
				cert.Properties.AuditComment = pointer.To(model.AuditComment)
			}

			if metadata.ResourceData.HasChanges("key_vault_certificate_id", "self_signed") {
				cert.Properties.CertificateSelfSigned = boolAsBooleanEnumCert(model.SelfSigned)
				cert.Properties.CertificateSignerResourceId = pointer.To(model.CertificateSignerID)
			}

			if err = client.CertificateObjectLocalRulestackCreateOrUpdateThenPoll(ctx, *id, cert); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			if err = rulestackClient.LocalRulestackscommitThenPoll(ctx, rulestackId); err != nil {
				return fmt.Errorf("committing Local RuleStack config for %s: %+v", id, err)
			}

			return nil
		},
	}
}
