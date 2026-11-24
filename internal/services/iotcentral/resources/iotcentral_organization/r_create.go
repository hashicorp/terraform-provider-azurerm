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
	"github.com/jackofallops/kermit/sdk/iotcentral/2022-10-31-preview/iotcentral"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTCentral
			var state IotCentralOrganizationModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			appId, err := apps.ParseIotAppID(state.IotCentralApplicationId)
			if err != nil {
				return err
			}

			// TODO: import check

			app, err := client.AppsClient.Get(ctx, *appId)
			if err != nil || app.Model == nil {
				return fmt.Errorf("checking for the presence of existing %q: %+v", appId, err)
			}

			orgClient, err := client.OrganizationsClient(ctx, *app.Model.Properties.Subdomain)
			if err != nil {
				return fmt.Errorf("creating organization client: %+v", err)
			}

			model := iotcentral.Organization{
				DisplayName: &state.DisplayName,
			}

			if state.ParentOrganizationId != "" {
				model.Parent = &state.ParentOrganizationId
			}

			org, err := orgClient.Create(ctx, state.OrganizationId, model)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", state.OrganizationId, err)
			}

			metadata.SetID(parse.NewOrganizationID(appId.SubscriptionId, appId.ResourceGroupName, appId.IotAppName, *org.ID))
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
