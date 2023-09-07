// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iotcentral

import (
	"context"
	"fmt"
	"time"

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
	SubDomain            string `tfschema:"sub_domain"`
	OrganizationId       string `tfschema:"organization_id"`
	DisplayName          string `tfschema:"display_name"`
	ParentOrganizationId string `tfschema:"parent_organization_id"`
}

func (r IotCentralOrganizationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"sub_domain": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"organization_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"display_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"parent_organization_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
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

			orgClient, err := client.OrganizationsClient(ctx, state.SubDomain)
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

			orgId, err := parse.NewOrganizationID(state.SubDomain, client.Endpoint.Name(), *org.ID)
			if err != nil {
				return err
			}

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
			id, err := parse.ParseOrganizationID(metadata.ResourceData.Id(), metadata.ResourceData.Get("sub_domain").(string), client.Endpoint.Name())
			if err != nil {
				return err
			}

			orgClient, err := client.OrganizationsClient(ctx, id.SubDomain)
			if err != nil {
				return fmt.Errorf("creating organization client: %+v", err)
			}

			org, err := orgClient.Get(ctx, id.OrganizationId)
			if err != nil {
				if org.ID == nil || *org.ID == "" {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := IotCentralOrganizationModel{
				SubDomain:      id.SubDomain,
				OrganizationId: id.OrganizationId,
				DisplayName:    *org.DisplayName,
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

			id, err := parse.ParseOrganizationID(metadata.ResourceData.Id(), state.SubDomain, client.Endpoint.Name())
			if err != nil {
				return err
			}

			orgClient, err := client.OrganizationsClient(ctx, id.SubDomain)
			if err != nil {
				return fmt.Errorf("creating organization client: %+v", err)
			}

			existing, err := orgClient.Get(ctx, id.OrganizationId)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
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

			id, err := parse.ParseOrganizationID(metadata.ResourceData.Id(), state.SubDomain, client.Endpoint.Name())
			if err != nil {
				return err
			}

			orgClient, err := client.OrganizationsClient(ctx, id.SubDomain)
			if err != nil {
				return fmt.Errorf("creating organization client: %+v", err)
			}

			_, err = orgClient.Remove(ctx, id.OrganizationId)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
