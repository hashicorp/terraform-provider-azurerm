// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package email_communication_service_domain

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2026-03-18/domains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
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

			if metadata.ResourceData.HasChange("user_engagement_tracking_enabled") {
				userEngagementTracking := domains.UserEngagementTrackingDisabled
				if model.UserEngagementTrackingEnabled {
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
