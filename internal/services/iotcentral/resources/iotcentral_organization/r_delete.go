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

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTCentral
			var state IotCentralOrganizationModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

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

			if _, err = orgClient.Remove(ctx, id.Name); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
