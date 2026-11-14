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
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.CertificateObjectLocalRulestackResources

			id, err := certificateobjectlocalrulestackresources.ParseLocalRulestackCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state LocalRuleStackCertificateModel

			existing, err := client.CertificateObjectLocalRulestackGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			state.Name = id.CertificateName
			state.RuleStackID = certificateobjectlocalrulestackresources.NewLocalRulestackID(id.SubscriptionId, id.ResourceGroupName, id.LocalRulestackName).ID()

			if model := existing.Model; model != nil {
				props := model.Properties

				state.AuditComment = pointer.From(props.AuditComment)
				state.CertificateSignerID = pointer.From(props.CertificateSignerResourceId)
				state.Description = pointer.From(props.Description)
				state.SelfSigned = boolEnumAsBoolCert(props.CertificateSelfSigned)
			}

			return metadata.Encode(&state)
		},
	}
}

func boolAsBooleanEnumCert(input bool) certificateobjectlocalrulestackresources.BooleanEnum {
	if input {
		return certificateobjectlocalrulestackresources.BooleanEnumTRUE
	}

	return certificateobjectlocalrulestackresources.BooleanEnumFALSE
}

func boolEnumAsBoolCert(input certificateobjectlocalrulestackresources.BooleanEnum) bool {
	return input == certificateobjectlocalrulestackresources.BooleanEnumTRUE
}
