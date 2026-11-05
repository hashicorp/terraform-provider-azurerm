// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cost_anomaly_alert

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2025-03-01/scheduledactionoperationgroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (Resource) Read() sdk.ResourceFunc {
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
