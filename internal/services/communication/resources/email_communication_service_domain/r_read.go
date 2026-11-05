// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package email_communication_service_domain

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2026-03-18/domains"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2026-03-18/emailservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/communication/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Communication.DomainClient

			state := EmailCommunicationServiceDomainResourceModel{}

			id, err := domains.ParseDomainID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state.Name = id.DomainName
			state.EMailServiceID = emailservices.NewEmailServiceID(id.SubscriptionId, id.ResourceGroupName, id.EmailServiceName).ID()

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.DomainManagement = string(props.DomainManagement)

					state.FromSenderDomain = pointer.From(props.FromSenderDomain)

					state.MailFromSenderDomain = pointer.From(props.MailFromSenderDomain)

					if props.UserEngagementTracking != nil {
						state.UserEngagementTrackingEnabled = *props.UserEngagementTracking == domains.UserEngagementTrackingEnabled
					}

					domainVerificationRecords := EmailDomainVerificationRecords{}

					if verificationRecords := props.VerificationRecords; verificationRecords != nil {
						if record := verificationRecords.DKIM; record != nil {
							domainVerificationRecords.DKIM = helper.DomainVerificationRecordsToModel(record)
						}

						if record := verificationRecords.DKIM2; record != nil {
							domainVerificationRecords.DKIM2 = helper.DomainVerificationRecordsToModel(record)
						}

						if record := verificationRecords.DMARC; record != nil {
							domainVerificationRecords.DMARC = helper.DomainVerificationRecordsToModel(record)
						}

						if record := verificationRecords.Domain; record != nil {
							domainVerificationRecords.Domain = helper.DomainVerificationRecordsToModel(record)
						}

						if record := verificationRecords.SPF; record != nil {
							domainVerificationRecords.SPF = helper.DomainVerificationRecordsToModel(record)
						}
					}

					state.VerificationRecords = []EmailDomainVerificationRecords{domainVerificationRecords}
				}

				state.Tags = pointer.From(model.Tags)
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}

			return metadata.Encode(&state)
		},
	}
}

type EmailDomainVerificationRecords struct {
	Domain []helper.DomainVerificationRecords `tfschema:"domain"`
	SPF    []helper.DomainVerificationRecords `tfschema:"spf"`
	DMARC  []helper.DomainVerificationRecords `tfschema:"dmarc"`
	DKIM   []helper.DomainVerificationRecords `tfschema:"dkim"`
	DKIM2  []helper.DomainVerificationRecords `tfschema:"dkim2"`
}
