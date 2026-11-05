// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package email_communication_service

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2026-03-18/emailservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Communication.EmailServicesClient

			var model EmailCommunicationServiceResourceModel

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id, err := emailservices.ParseEmailServiceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil || existing.Model == nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			emailService := *existing.Model

			props := pointer.From(emailService.Properties)

			if metadata.ResourceData.HasChange("data_location") {
				props.DataLocation = model.DataLocation
			}

			if metadata.ResourceData.HasChange("tags") {
				emailService.Tags = pointer.To(model.Tags)
			}

			emailService.Properties = &props

			if err := client.CreateOrUpdateThenPoll(ctx, *id, emailService); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
