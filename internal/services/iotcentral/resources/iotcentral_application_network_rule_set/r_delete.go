// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iotcentral_application_network_rule_set

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/iotcentral/2021-11-01-preview/apps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTCentral.AppsClient
			id, err := apps.ParseIotAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			app, err := client.Get(ctx, *id)
			if err != nil || app.Model == nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if app.Model.Properties == nil {
				app.Model.Properties = &apps.AppProperties{}
			}

			app.Model.Properties.NetworkRuleSets = nil

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *app.Model); err != nil {
				return fmt.Errorf("deleting Network Rule Set of %s: %+v", id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
