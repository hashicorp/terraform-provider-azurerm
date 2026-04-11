// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apigateway"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apigatewayconfigconnection"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ApiManagementGatewayConfigConnectionModel struct {
	Name                   string   `tfschema:"name"`
	ApiManagementGatewayId string   `tfschema:"api_management_gateway_id"`
	WorkspaceId            string   `tfschema:"workspace_id"`
	Hostnames              []string `tfschema:"hostnames"`
}

type ApiManagementGatewayWorkspaceConnectionResource struct{}

var _ sdk.Resource = ApiManagementGatewayWorkspaceConnectionResource{}

func (r ApiManagementGatewayWorkspaceConnectionResource) ResourceType() string {
	return "azurerm_api_management_gateway_workspace_connection"
}

func (r ApiManagementGatewayWorkspaceConnectionResource) ModelObject() interface{} {
	return &ApiManagementGatewayConfigConnectionModel{}
}

func (r ApiManagementGatewayWorkspaceConnectionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return apigatewayconfigconnection.ValidateConfigConnectionID
}

func (r ApiManagementGatewayWorkspaceConnectionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z](?:[a-zA-Z0-9-]{0,28}[a-zA-Z0-9])?$"),
				"The `name` must be between 1 and 30 characters long and can only include letters, numbers, and hyphens. The first character must be a letter and last character must be a letter or a number.",
			),
		},

		"api_management_gateway_id": commonschema.ResourceIDReferenceRequiredForceNew(&apigateway.GatewayId{}),

		"workspace_id": commonschema.ResourceIDReferenceRequiredForceNew(&workspace.WorkspaceId{}),

		"hostnames": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func (r ApiManagementGatewayWorkspaceConnectionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ApiManagementGatewayWorkspaceConnectionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.ApiGatewayConfigConnectionClient

			var model ApiManagementGatewayConfigConnectionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			gatewayId, err := apigateway.ParseGatewayID(model.ApiManagementGatewayId)
			if err != nil {
				return err
			}

			id := apigatewayconfigconnection.NewConfigConnectionID(gatewayId.SubscriptionId, gatewayId.ResourceGroupName, gatewayId.GatewayName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := apigatewayconfigconnection.ApiManagementGatewayConfigConnectionResource{
				Properties: apigatewayconfigconnection.GatewayConfigConnectionBaseProperties{
					SourceId: pointer.To(model.WorkspaceId),
				},
			}

			if len(model.Hostnames) > 0 {
				payload.Properties.Hostnames = pointer.To(model.Hostnames)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ApiManagementGatewayWorkspaceConnectionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.ApiGatewayConfigConnectionClient

			id, err := apigatewayconfigconnection.ParseConfigConnectionID(metadata.ResourceData.Id())
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

			state := ApiManagementGatewayConfigConnectionModel{
				Name:                   id.ConfigConnectionName,
				ApiManagementGatewayId: apigateway.NewGatewayID(id.SubscriptionId, id.ResourceGroupName, id.GatewayName).ID(),
			}

			if model := resp.Model; model != nil {
				state.WorkspaceId = pointer.From(model.Properties.SourceId)
				state.Hostnames = pointer.From(model.Properties.Hostnames)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ApiManagementGatewayWorkspaceConnectionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.ApiGatewayConfigConnectionClient

			id, err := apigatewayconfigconnection.ParseConfigConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, apigatewayconfigconnection.DefaultDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
