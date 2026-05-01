// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package costmanagement

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2023-08-01/scheduledactions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2023-08-01/views"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = AnomalyAlertResource{}

type AnomalyAlertResource struct{}

type AnomalyAlertModel struct {
	Name              string   `tfschema:"name"`
	DisplayName       string   `tfschema:"display_name"`
	SubscriptionId    string   `tfschema:"subscription_id"`
	NotificationEmail string   `tfschema:"notification_email"`
	EmailSubject      string   `tfschema:"email_subject"`
	EmailAddresses    []string `tfschema:"email_addresses"`
	Message           string   `tfschema:"message"`
}

func (AnomalyAlertResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CostAnomalyAlertName,
			// action names can contain only alphanumeric characters and hyphens.
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			// adding 2026/04 api limitations behind 5.0 flag incase existing resources are still allowed with previous limits
			ValidateFunc: func() pluginsdk.SchemaValidateFunc {
				if features.FivePointOh() {
					return validation.StringLenBetween(1, 25)
				}
				return validation.StringIsNotEmpty
			}(),
		},

		"subscription_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Computed:     true,
			ValidateFunc: commonids.ValidateSubscriptionID,
		},

		"notification_email": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"email_subject": {
			Type:     pluginsdk.TypeString,
			Required: true,
			// adding 2026/04 api limitations behind 5.0 flag incase existing resources are still allowed with previous limits
			ValidateFunc: func() pluginsdk.SchemaValidateFunc {
				if features.FivePointOh() {
					return validation.StringLenBetween(1, 50)
				}
				return validation.StringLenBetween(1, 70)
			}(),
		},

		"email_addresses": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"message": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 250),
		},
	}
}

func (AnomalyAlertResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (AnomalyAlertResource) ModelObject() interface{} {
	return &AnomalyAlertModel{}
}

func (AnomalyAlertResource) ResourceType() string {
	return "azurerm_cost_anomaly_alert"
}

func (AnomalyAlertResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return scheduledactions.ValidateScopedScheduledActionID
}

func (r AnomalyAlertResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CostManagement.ScheduledActionsClient

			var config AnomalyAlertModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := config.SubscriptionId
			if subscriptionId == "" {
				subscriptionId = commonids.NewSubscriptionID(metadata.Client.Account.SubscriptionId).ID()
			}
			id := scheduledactions.NewScopedScheduledActionID(subscriptionId, config.Name)

			existing, err := client.GetByScope(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			viewId := views.NewScopedViewID(subscriptionId, "ms:DailyAnomalyByResourceGroup")

			schedule := scheduledactions.ScheduleProperties{
				Frequency: scheduledactions.ScheduleFrequencyDaily,
			}
			schedule.SetEndDateAsTime(time.Now().AddDate(1, 0, 0))
			schedule.SetStartDateAsTime(time.Now())

			notificationEmail := config.EmailAddresses[0]
			if config.NotificationEmail != "" {
				notificationEmail = config.NotificationEmail
			}

			param := scheduledactions.ScheduledAction{
				Kind: pointer.To(scheduledactions.ScheduledActionKindInsightAlert),
				Properties: &scheduledactions.ScheduledActionProperties{
					DisplayName: config.DisplayName,
					Status:      scheduledactions.ScheduledActionStatusEnabled,
					ViewId:      viewId.ID(),
					FileDestination: &scheduledactions.FileDestination{
						FileFormats: &[]scheduledactions.FileFormat{},
					},
					NotificationEmail: &notificationEmail,
					Notification: scheduledactions.NotificationProperties{
						Subject: config.EmailSubject,
						Message: pointer.To(config.Message),
						To:      config.EmailAddresses,
					},
					Schedule: schedule,
				},
			}
			if _, err := client.CreateOrUpdateByScope(ctx, id, param, scheduledactions.DefaultCreateOrUpdateByScopeOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r AnomalyAlertResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CostManagement.ScheduledActionsClient

			id, err := scheduledactions.ParseScopedScheduledActionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config AnomalyAlertModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.GetByScope(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				if model.ETag == nil {
					return fmt.Errorf("add %s: etag was nil", *id)
				}
			}

			subscriptionId := config.SubscriptionId
			if subscriptionId == "" {
				subscriptionId = commonids.NewSubscriptionID(metadata.Client.Account.SubscriptionId).ID()
			}
			viewId := views.NewScopedViewID(subscriptionId, "ms:DailyAnomalyByResourceGroup")

			schedule := scheduledactions.ScheduleProperties{
				Frequency: scheduledactions.ScheduleFrequencyDaily,
			}
			schedule.SetEndDateAsTime(time.Now().AddDate(1, 0, 0))
			schedule.SetStartDateAsTime(time.Now())

			notificationEmail := config.EmailAddresses[0]
			if config.NotificationEmail != "" {
				notificationEmail = config.NotificationEmail
			}

			param := scheduledactions.ScheduledAction{
				Kind: pointer.To(scheduledactions.ScheduledActionKindInsightAlert),
				ETag: resp.Model.ETag,
				Properties: &scheduledactions.ScheduledActionProperties{
					DisplayName:       config.DisplayName,
					Status:            scheduledactions.ScheduledActionStatusEnabled,
					ViewId:            viewId.ID(),
					NotificationEmail: &notificationEmail,
					Notification: scheduledactions.NotificationProperties{
						Subject: config.EmailSubject,
						Message: pointer.To(config.Message),
						To:      config.EmailAddresses,
					},
					Schedule: schedule,
				},
			}
			if _, err := client.CreateOrUpdateByScope(ctx, *id, param, scheduledactions.DefaultCreateOrUpdateByScopeOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (AnomalyAlertResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CostManagement.ScheduledActionsClient

			id, err := scheduledactions.ParseScopedScheduledActionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.GetByScope(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := AnomalyAlertModel{}

			if model := resp.Model; model != nil {
				state.Name = pointer.From(model.Name)
				if props := model.Properties; props != nil {
					state.DisplayName = props.DisplayName
					if props.Scope != nil {
						state.SubscriptionId = fmt.Sprint("/", *props.Scope)
					}
					state.EmailSubject = props.Notification.Subject
					state.NotificationEmail = pointer.From(props.NotificationEmail)
					state.EmailAddresses = props.Notification.To
					state.Message = pointer.From(props.Notification.Message)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (AnomalyAlertResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CostManagement.ScheduledActionsClient

			id, err := scheduledactions.ParseScopedScheduledActionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err = client.DeleteByScope(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
