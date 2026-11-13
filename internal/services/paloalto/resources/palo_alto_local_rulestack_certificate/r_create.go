// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_local_rulestack_certificate

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/certificateobjectlocalrulestackresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/localrulestackresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.CertificateObjectLocalRulestackResources
			rulestackClient := metadata.Client.PaloAlto.LocalRulestackResources

			model := LocalRuleStackCertificateModel{}
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			rulestackId, err := localrulestackresources.ParseLocalRulestackID(model.RuleStackID)
			if err != nil {
				return err
			}

			locks.ByID(rulestackId.ID())
			defer locks.UnlockByID(rulestackId.ID())

			id := certificateobjectlocalrulestackresources.NewLocalRulestackCertificateID(rulestackId.SubscriptionId, rulestackId.ResourceGroupName, rulestackId.LocalRulestackName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.CertificateObjectLocalRulestackGet(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			props := certificateobjectlocalrulestackresources.CertificateObject{
				CertificateSelfSigned: boolAsBooleanEnumCert(model.SelfSigned),
			}

			if model.AuditComment != "" {
				props.AuditComment = pointer.To(model.AuditComment)
			}

			if model.CertificateSignerID != "" {
				props.CertificateSignerResourceId = pointer.To(model.CertificateSignerID)
			}

			if model.Description != "" {
				props.Description = pointer.To(model.Description)
			}

			cert := certificateobjectlocalrulestackresources.CertificateObjectLocalRulestackResource{
				Properties: props,
			}

			if _, err = client.CertificateObjectLocalRulestackCreateOrUpdate(ctx, id, cert); err != nil {
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
