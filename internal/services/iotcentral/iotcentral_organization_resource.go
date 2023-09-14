// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iotcentral

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/iotcentral/2021-11-01-preview/apps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iotcentral/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iotcentral/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	dataplane "github.com/tombuildsstuff/kermit/sdk/iotcentral/2022-10-31-preview/iotcentral"
)

type IotCentralOrganizationResource struct{}

var (
	_ sdk.ResourceWithUpdate = IotCentralOrganizationResource{}
)

type IotCentralOrganizationModel struct {
	IotCentralApplicationId string `tfschema:"iotcentral_application_id"`
	OrganizationId          string `tfschema:"organization_id"`
	DisplayName             string `tfschema:"display_name"`
	ParentOrganizationId    string `tfschema:"parent_organization_id"`
}

func (r IotCentralOrganizationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"iotcentral_application_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: apps.ValidateIotAppID,
		},
		"organization_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.OrganizationOrganizationID,
		},
		"display_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"parent_organization_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validate.OrganizationOrganizationID,
		},
	}
}

func (r IotCentralOrganizationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r IotCentralOrganizationResource) ResourceType() string {
	return "azurerm_iotcentral_organization"
}

func (r IotCentralOrganizationResource) ModelObject() interface{} {
	return &IotCentralOrganizationModel{}
}

func (r IotCentralOrganizationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.OrganizationID
}

func (r IotCentralOrganizationResource) Create() sdk.ResourceFunc {
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

			app, err := client.AppsClient.Get(ctx, *appId)
			if err != nil || app.Model == nil {
				return fmt.Errorf("checking for the presence of existing %q: %+v", appId, err)
			}

			orgClient, err := client.OrganizationsClient(ctx, *app.Model.Properties.Subdomain)
			if err != nil {
				return fmt.Errorf("creating organization client: %+v", err)
			}

			model := dataplane.Organization{
				DisplayName: &state.DisplayName,
			}

			if state.ParentOrganizationId != "" {
				model.Parent = &state.ParentOrganizationId
			}

			org, err := orgClient.Create(ctx, state.OrganizationId, model)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", state.OrganizationId, err)
			}

			orgId := parse.NewOrganizationID(appId.SubscriptionId, appId.ResourceGroupName, appId.IotAppName, *org.ID)

			metadata.SetID(orgId)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r IotCentralOrganizationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTCentral
			id, err := parse.OrganizationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			appId := apps.NewIotAppID(id.SubscriptionId, id.ResourceGroup, id.IotAppName)
			if err != nil {
				return err
			}

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

func (r IotCentralOrganizationResource) Update() sdk.ResourceFunc {
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
			if err != nil {
				return err
			}

			app, err := client.AppsClient.Get(ctx, appId)
			if err != nil || app.Model == nil {
				return metadata.MarkAsGone(id)
			}

			orgClient, err := client.OrganizationsClient(ctx, *app.Model.Properties.Subdomain)
			if err != nil {
				return fmt.Errorf("creating organization client: %+v", err)
			}

			existing, err := orgClient.Get(ctx, id.Name)
			if err != nil {
				if existing.ID == nil || *existing.ID == "" {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if metadata.ResourceData.HasChange("display_name") {
				existing.DisplayName = &state.DisplayName
			}

			_, err = orgClient.Update(ctx, *existing.ID, existing, "*")
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r IotCentralOrganizationResource) Delete() sdk.ResourceFunc {
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
			if err != nil {
				return err
			}

			app, err := client.AppsClient.Get(ctx, appId)
			if err != nil || app.Model == nil {
				return metadata.MarkAsGone(id)
			}

			orgClient, err := client.OrganizationsClient(ctx, *app.Model.Properties.Subdomain)
			if err != nil {
				return fmt.Errorf("creating organization client: %+v", err)
			}

			_, err = orgClient.Remove(ctx, id.Name)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
