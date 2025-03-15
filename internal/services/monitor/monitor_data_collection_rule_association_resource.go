// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-11/datacollectionendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-11/datacollectionruleassociations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-11/datacollectionrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DataCollectionRuleAssociationModel struct {
	Name                     string `tfschema:"name"`
	TargetResourceId         string `tfschema:"target_resource_id"`
	DataCollectionEndpointId string `tfschema:"data_collection_endpoint_id"`
	DataCollectionRuleId     string `tfschema:"data_collection_rule_id"`
	Description              string `tfschema:"description"`
}

type DataCollectionRuleAssociationResource struct{}

func (r DataCollectionRuleAssociationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"target_resource_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"name": {
			// TODO: should this be hard-coded in the Create?
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      "configurationAccessEndpoint",
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"data_collection_endpoint_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: datacollectionendpoints.ValidateDataCollectionEndpointID,
			ExactlyOneOf: []string{"data_collection_endpoint_id", "data_collection_rule_id"},
		},

		"data_collection_rule_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: datacollectionrules.ValidateDataCollectionRuleID,
			ExactlyOneOf: []string{"data_collection_endpoint_id", "data_collection_rule_id"},
			RequiredWith: []string{"name"},
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}

func (r DataCollectionRuleAssociationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DataCollectionRuleAssociationResource) ResourceType() string {
	return "azurerm_monitor_data_collection_rule_association"
}

func (r DataCollectionRuleAssociationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return datacollectionruleassociations.ValidateScopedDataCollectionRuleAssociationID
}

func (r DataCollectionRuleAssociationResource) ModelObject() interface{} {
	return &DataCollectionRuleAssociationModel{}
}

func (r DataCollectionRuleAssociationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			metadata.Logger.Info("Decoding state..")
			var model DataCollectionRuleAssociationModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := metadata.Client.Monitor.DataCollectionRuleAssociationsClient

			id := datacollectionruleassociations.NewScopedDataCollectionRuleAssociationID(model.TargetResourceId, model.Name)
			metadata.Logger.Infof("creating %s", id)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			input := datacollectionruleassociations.DataCollectionRuleAssociationProxyOnlyResource{
				Name: utils.String(model.Name),
				Properties: &datacollectionruleassociations.DataCollectionRuleAssociation{
					Description: utils.String(model.Description),
				},
			}

			if model.DataCollectionEndpointId != "" {
				input.Properties.DataCollectionEndpointId = utils.String(model.DataCollectionEndpointId)
			}
			if model.DataCollectionRuleId != "" {
				input.Properties.DataCollectionRuleId = utils.String(model.DataCollectionRuleId)
			}

			if _, err := client.Create(ctx, id, input); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r DataCollectionRuleAssociationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.DataCollectionRuleAssociationsClient
			id, err := datacollectionruleassociations.ParseScopedDataCollectionRuleAssociationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("retrieving %s", *id)
			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					metadata.Logger.Infof("%s was not found - removing from state!", *id)
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var description, dataCollectionEndpointId, dataCollectionRuleId string

			if model := resp.Model; model != nil {
				if prop := model.Properties; prop != nil {
					dataCollectionEndpointId = flattenStringPtr(prop.DataCollectionEndpointId)
					dataCollectionRuleId = flattenStringPtr(prop.DataCollectionRuleId)
					description = flattenStringPtr(prop.Description)
				}
			}

			return metadata.Encode(&DataCollectionRuleAssociationModel{
				Name:                     id.DataCollectionRuleAssociationName,
				TargetResourceId:         id.ResourceUri,
				DataCollectionEndpointId: dataCollectionEndpointId,
				DataCollectionRuleId:     dataCollectionRuleId,
				Description:              description,
			})
		},
		Timeout: 5 * time.Minute,
	}
}

func (r DataCollectionRuleAssociationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := datacollectionruleassociations.ParseScopedDataCollectionRuleAssociationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("updating %s..", *id)
			client := metadata.Client.Monitor.DataCollectionRuleAssociationsClient
			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if resp.Model == nil {
				return fmt.Errorf("unexpected null model of %s", *id)
			}
			existing := resp.Model
			if existing.Properties == nil {
				return fmt.Errorf("unexpected null properties of %s", *id)
			}

			var model DataCollectionRuleAssociationModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			if metadata.ResourceData.HasChange("data_collection_endpoint_id") {
				if model.DataCollectionEndpointId != "" {
					existing.Properties.DataCollectionEndpointId = utils.String(model.DataCollectionEndpointId)
				} else {
					existing.Properties.DataCollectionEndpointId = nil
				}
			}

			if metadata.ResourceData.HasChange("data_collection_rule_id") {
				if model.DataCollectionRuleId != "" {
					existing.Properties.DataCollectionRuleId = utils.String(model.DataCollectionRuleId)
				} else {
					existing.Properties.DataCollectionRuleId = nil
				}
			}

			if metadata.ResourceData.HasChange("description") {
				existing.Properties.Description = utils.String(model.Description)
			}

			if _, err := client.Create(ctx, *id, *existing); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r DataCollectionRuleAssociationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.DataCollectionRuleAssociationsClient
			id, err := datacollectionruleassociations.ParseScopedDataCollectionRuleAssociationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s..", *id)
			resp, err := client.Delete(ctx, *id)
			if err != nil && !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
