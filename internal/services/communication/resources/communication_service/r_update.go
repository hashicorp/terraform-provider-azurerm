// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package communication_service

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2026-03-18/communicationservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Communication.ServiceClient

			var model CommunicationServiceResourceModel

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id, err := communicationservices.ParseCommunicationServiceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil || existing.Model == nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			commService := *existing.Model

			props := pointer.From(commService.Properties)

			if metadata.ResourceData.HasChange("data_location") {
				props.DataLocation = model.DataLocation
			}

			existing.Model.Properties = &props

			if metadata.ResourceData.HasChange("tags") {
				commService.Tags = pointer.To(model.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, commService); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
