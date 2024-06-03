// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package communication

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2023-03-31/domains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/communication/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/communication/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = EmailCommunicationServiceDomainResource{}

type EmailCommunicationServiceDomainResource struct{}

type EmailCommunicationServiceDomainResourceModel struct {
	Name                   string            `tfschema:"name"`
	ResourceGroupName      string            `tfschema:"resource_group_name"`
	EMailServiceName       string            `tfschema:"email_service_name"`
	DomainManagement       string            `tfschema:"domain_management"`
	UserEngagementTracking bool              `tfschema:"user_engagement_tracking"`
	Tags                   map[string]string `tfschema:"tags"`

	FromSenderDomain     string                           `tfschema:"from_sender_domain"`
	MailFromSenderDomain string                           `tfschema:"mail_from_sender_domain"`
	VerificationRecords  []EmailDomainVerificationRecords `tfschema:"verification_records"`
}

type EmailDomainVerificationRecords struct {
	Domain []helper.DomainVerificationRecords `tfschema:"domain"`
	SPF    []helper.DomainVerificationRecords `tfschema:"spf"`
	DMARC  []helper.DomainVerificationRecords `tfschema:"dmarc"`
	DKIM   []helper.DomainVerificationRecords `tfschema:"dkim"`
	DKIM2  []helper.DomainVerificationRecords `tfschema:"dkim2"`
}

func (EmailCommunicationServiceDomainResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"email_service_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CommunicationServiceName,
		},

		"domain_management": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(domains.PossibleValuesForDomainManagement(), false),
		},

		"user_engagement_tracking": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
		},

		"tags": commonschema.Tags(),
	}
}

func (EmailCommunicationServiceDomainResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"from_sender_domain": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"mail_from_sender_domain": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"verification_records": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"domain": helper.DomainVerificationRecordsCommonSchema(),
					"dkim":   helper.DomainVerificationRecordsCommonSchema(),
					"dkim2":  helper.DomainVerificationRecordsCommonSchema(),
					"dmarc":  helper.DomainVerificationRecordsCommonSchema(),
					"spf":    helper.DomainVerificationRecordsCommonSchema(),
				},
			},
		},
	}
}

func (EmailCommunicationServiceDomainResource) ModelObject() interface{} {
	return &EmailCommunicationServiceDomainResourceModel{}
}

func (EmailCommunicationServiceDomainResource) ResourceType() string {
	return "azurerm_email_communication_service_domain"
}

func (r EmailCommunicationServiceDomainResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.Communication.DomainClient

			var model EmailCommunicationServiceDomainResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := domains.NewDomainID(subscriptionId, model.ResourceGroupName, model.EMailServiceName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &domains.DomainProperties{
				DomainManagement: domains.DomainManagement(model.DomainManagement),
			}

			properties.UserEngagementTracking = pointer.To(domains.UserEngagementTrackingDisabled)
			if model.UserEngagementTracking {
				properties.UserEngagementTracking = pointer.To(domains.UserEngagementTrackingEnabled)
			}

			param := domains.DomainResource{
				// The location is always `global` from the Azure Portal
				Location:   location.Normalize("global"),
				Properties: properties,
				Tags:       pointer.To(model.Tags),
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r EmailCommunicationServiceDomainResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Communication.DomainClient

			var model EmailCommunicationServiceDomainResourceModel

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id, err := domains.ParseDomainID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil || existing.Model == nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			domain := *existing.Model

			props := pointer.From(domain.Properties)

			if metadata.ResourceData.HasChange("user_engagement_tracking") {
				userEngagementTracking := domains.UserEngagementTrackingDisabled
				if model.UserEngagementTracking {
					userEngagementTracking = domains.UserEngagementTrackingEnabled
				}

				props.UserEngagementTracking = pointer.To(userEngagementTracking)
			}

			if metadata.ResourceData.HasChange("tags") {
				domain.Tags = pointer.To(model.Tags)
			}

			domain.Properties = &props

			if err := client.CreateOrUpdateThenPoll(ctx, *id, domain); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (EmailCommunicationServiceDomainResource) Read() sdk.ResourceFunc {
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
			state.ResourceGroupName = id.ResourceGroupName
			state.EMailServiceName = id.EmailServiceName

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.DomainManagement = string(props.DomainManagement)

					if props.FromSenderDomain != nil {
						state.FromSenderDomain = pointer.From(props.FromSenderDomain)
					}

					if props.MailFromSenderDomain != nil {
						state.MailFromSenderDomain = pointer.From(props.MailFromSenderDomain)
					}

					if props.UserEngagementTracking != nil {
						state.UserEngagementTracking = *props.UserEngagementTracking == domains.UserEngagementTrackingEnabled

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

			return metadata.Encode(&state)
		},
	}
}

func (EmailCommunicationServiceDomainResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Communication.DomainClient

			id, err := domains.ParseDomainID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (EmailCommunicationServiceDomainResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return domains.ValidateDomainID
}
