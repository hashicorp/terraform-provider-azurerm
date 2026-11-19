// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iotcentral_organization

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/iotcentral/2021-11-01-preview/apps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iotcentral/parse"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTCentral
			id, err := parse.OrganizationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			appId := apps.NewIotAppID(id.SubscriptionId, id.ResourceGroup, id.IotAppName)

			app, err := client.AppsClient.Get(ctx, appId)
			if err != nil || app.Model == nil {
				return metadata.MarkAsGone(id)
			}

			orgClient, err := client.OrganizationsClient(ctx, *app.Model.Properties.Subdomain)
			if err != nil {
				return fmt.Errorf("creating organization client: %+v", err)
			}

			org, err := orgClient.Get(ctx, id.Name)
			if err != nil {
				if org.ID == nil || *org.ID == "" {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := IotCentralOrganizationModel{
				IotCentralApplicationId: appId.ID(),
				OrganizationId:          id.Name,
				DisplayName:             *org.DisplayName,
			}

			if org.Parent != nil {
				state.ParentOrganizationId = *org.Parent
			}

			return metadata.Encode(&state)
		},
		Timeout: 5 * time.Minute,
	}
}
