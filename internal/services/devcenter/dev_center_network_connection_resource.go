// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devcenter

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/networkconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devcenter/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.Resource           = DevCenterNetworkConnectionResource{}
	_ sdk.ResourceWithUpdate = DevCenterNetworkConnectionResource{}
)

type DevCenterNetworkConnectionResource struct{}

func (r DevCenterNetworkConnectionResource) ModelObject() interface{} {
	return &DevCenterNetworkConnectionResourceModel{}
}

type DevCenterNetworkConnectionResourceModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	DomainJoinType    string            `tfschema:"domain_join_type"`
	SubnetId          string            `tfschema:"subnet_id"`
	DomainName        string            `tfschema:"domain_name"`
	DomainPassword    string            `tfschema:"domain_password"`
	DomainUsername    string            `tfschema:"domain_username"`
	OrganizationUnit  string            `tfschema:"organization_unit"`
	Tags              map[string]string `tfschema:"tags"`
}

func (r DevCenterNetworkConnectionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return networkconnections.ValidateNetworkConnectionID
}

func (r DevCenterNetworkConnectionResource) ResourceType() string {
	return "azurerm_dev_center_network_connection"
}

func (r DevCenterNetworkConnectionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DevCenterNetworkConnectionName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"domain_join_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(networkconnections.PossibleValuesForDomainJoinType(), false),
		},

		"subnet_id": commonschema.ResourceIDReferenceRequired(&commonids.SubnetId{}),

		"domain_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.DevCenterNetworkConnectionDomainName,
		},

		"domain_password": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"domain_username": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.DevCenterNetworkConnectionDomainUsername,
		},

		"organization_unit": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.Tags(),
	}
}

func (r DevCenterNetworkConnectionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DevCenterNetworkConnectionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.NetworkConnections
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model DevCenterNetworkConnectionResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := networkconnections.NewNetworkConnectionID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := networkconnections.NetworkConnection{
				Location: location.Normalize(model.Location),
				Properties: &networkconnections.NetworkProperties{
					DomainJoinType: networkconnections.DomainJoinType(model.DomainJoinType),
					SubnetId:       pointer.To(model.SubnetId),
				},
				Tags: pointer.To(model.Tags),
			}

			if v := model.DomainName; v != "" {
				parameters.Properties.DomainName = pointer.To(v)
			}

			if v := model.DomainPassword; v != "" {
				parameters.Properties.DomainPassword = pointer.To(v)
			}

			if v := model.DomainUsername; v != "" {
				parameters.Properties.DomainUsername = pointer.To(v)
			}

			if v := model.OrganizationUnit; v != "" {
				parameters.Properties.OrganizationUnit = pointer.To(v)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r DevCenterNetworkConnectionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.NetworkConnections

			id, err := networkconnections.ParseNetworkConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := DevCenterNetworkConnectionResourceModel{
				Name:              id.NetworkConnectionName,
				ResourceGroupName: id.ResourceGroupName,
				DomainPassword:    metadata.ResourceData.Get("domain_password").(string),
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				if props := model.Properties; props != nil {
					state.SubnetId = pointer.From(props.SubnetId)
					state.DomainName = pointer.From(props.DomainName)
					state.DomainUsername = pointer.From(props.DomainUsername)
					state.OrganizationUnit = pointer.From(props.OrganizationUnit)

					if v := props.DomainJoinType; v != "" {
						state.DomainJoinType = string(v)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r DevCenterNetworkConnectionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.NetworkConnections

			id, err := networkconnections.ParseNetworkConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r DevCenterNetworkConnectionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.NetworkConnections

			id, err := networkconnections.ParseNetworkConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model DevCenterNetworkConnectionResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := networkconnections.NetworkConnectionUpdate{
				Properties: &networkconnections.NetworkConnectionUpdateProperties{},
			}

			if metadata.ResourceData.HasChange("subnet_id") {
				parameters.Properties.SubnetId = pointer.To(model.SubnetId)
			}

			if metadata.ResourceData.HasChange("domain_name") {
				parameters.Properties.DomainName = pointer.To(model.DomainName)
			}

			if metadata.ResourceData.HasChange("domain_password") {
				parameters.Properties.DomainPassword = pointer.To(model.DomainPassword)
			}

			if metadata.ResourceData.HasChange("domain_username") {
				parameters.Properties.DomainUsername = pointer.To(model.DomainUsername)
			}

			if metadata.ResourceData.HasChange("organization_unit") {
				parameters.Properties.OrganizationUnit = pointer.To(model.OrganizationUnit)
			}

			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = pointer.To(model.Tags)
			}

			if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
