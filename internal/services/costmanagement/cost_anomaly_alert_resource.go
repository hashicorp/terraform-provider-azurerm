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
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2025-03-01/scheduledactionoperationgroup"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2025-03-01/viewoperationgroup"
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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(1, 25),
		},

		"subscription_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: commonids.ValidateSubscriptionID,
		},

		"notification_email": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"email_subject": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(1, 50),
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

func (AnomalyAlertResource) ModelObject() any {
	return &AnomalyAlertModel{}
}

func (AnomalyAlertResource) ResourceType() string {
	return "azurerm_cost_anomaly_alert"
}

func (AnomalyAlertResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return scheduledactionoperationgroup.ValidateScopedScheduledActionID
}

func (r AnomalyAlertResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CostManagement.ScheduledActionOperationGroupClient

			var config AnomalyAlertModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := config.SubscriptionId
			if subscriptionId == "" {
				subscriptionId = commonids.NewSubscriptionID(metadata.Client.Account.SubscriptionId).ID()
			}
			id := scheduledactionoperationgroup.NewScopedScheduledActionID(subscriptionId, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.ScheduledActionsGetByScope(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			viewId := viewoperationgroup.NewScopedViewID(subscriptionId, "ms:DailyAnomalyByResourceGroup")

			schedule := scheduledactionoperationgroup.ScheduleProperties{
				Frequency: scheduledactionoperationgroup.ScheduleFrequencyDaily,
			}
			schedule.SetEndDateAsTime(time.Now().AddDate(1, 0, 0))
			schedule.SetStartDateAsTime(time.Now())

			notificationEmail := config.EmailAddresses[0]
			if config.NotificationEmail != "" {
				notificationEmail = config.NotificationEmail
			}

			param := scheduledactionoperationgroup.ScheduledAction{
				Kind: pointer.To(scheduledactionoperationgroup.ScheduledActionKindInsightAlert),
				Properties: &scheduledactionoperationgroup.ScheduledActionProperties{
					DisplayName: config.DisplayName,
					Status:      scheduledactionoperationgroup.ScheduledActionStatusEnabled,
					ViewId:      viewId.ID(),
					FileDestination: &scheduledactionoperationgroup.FileDestination{
						FileFormats: &[]scheduledactionoperationgroup.FileFormat{},
					},
					NotificationEmail: &notificationEmail,
					Notification: scheduledactionoperationgroup.NotificationProperties{
						Subject: config.EmailSubject,
						Message: pointer.To(config.Message),
						To:      config.EmailAddresses,
					},
					Schedule: schedule,
				},
			}
			if _, err := client.ScheduledActionsCreateOrUpdateByScope(ctx, id, param, scheduledactionoperationgroup.DefaultScheduledActionsCreateOrUpdateByScopeOperationOptions()); err != nil {
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
			client := metadata.Client.CostManagement.ScheduledActionOperationGroupClient

			id, err := scheduledactionoperationgroup.ParseScopedScheduledActionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config AnomalyAlertModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.ScheduledActionsGetByScope(ctx, *id)
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
			viewId := viewoperationgroup.NewScopedViewID(subscriptionId, "ms:DailyAnomalyByResourceGroup")

			schedule := scheduledactionoperationgroup.ScheduleProperties{
				Frequency: scheduledactionoperationgroup.ScheduleFrequencyDaily,
			}
			schedule.SetEndDateAsTime(time.Now().AddDate(1, 0, 0))
			schedule.SetStartDateAsTime(time.Now())

			notificationEmail := config.EmailAddresses[0]
			if config.NotificationEmail != "" {
				notificationEmail = config.NotificationEmail
			}

			param := scheduledactionoperationgroup.ScheduledAction{
				Kind: pointer.To(scheduledactionoperationgroup.ScheduledActionKindInsightAlert),
				ETag: resp.Model.ETag,
				Properties: &scheduledactionoperationgroup.ScheduledActionProperties{
					DisplayName:       config.DisplayName,
					Status:            scheduledactionoperationgroup.ScheduledActionStatusEnabled,
					ViewId:            viewId.ID(),
					NotificationEmail: &notificationEmail,
					Notification: scheduledactionoperationgroup.NotificationProperties{
						Subject: config.EmailSubject,
						Message: pointer.To(config.Message),
						To:      config.EmailAddresses,
					},
					Schedule: schedule,
				},
			}
			if _, err := client.ScheduledActionsCreateOrUpdateByScope(ctx, *id, param, scheduledactionoperationgroup.DefaultScheduledActionsCreateOrUpdateByScopeOperationOptions()); err != nil {
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
			client := metadata.Client.CostManagement.ScheduledActionOperationGroupClient

			id, err := scheduledactionoperationgroup.ParseScopedScheduledActionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.ScheduledActionsGetByScope(ctx, *id)
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
			client := metadata.Client.CostManagement.ScheduledActionOperationGroupClient

			id, err := scheduledactionoperationgroup.ParseScopedScheduledActionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err = client.ScheduledActionsDeleteByScope(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
