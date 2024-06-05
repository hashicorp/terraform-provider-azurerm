// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package advisor

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/advisor/2023-01-01/suppressions"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/advisor/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = AdvisorSuppressionResource{}

type AdvisorSuppressionResource struct{}

type AdvisorSuppressionResourceModel struct {
	Name             string `tfschema:"name"`
	SuppressionID    string `tfschema:"suppression_id"`
	RecommendationID string `tfschema:"recommendation_id"`
	ResourceID       string `tfschema:"resource_id"`
	TTL              string `tfschema:"ttl"`
}

func (AdvisorSuppressionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"recommendation_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"resource_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},
		"ttl": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validate.Duration,
		},
	}
}

func (AdvisorSuppressionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"suppression_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (AdvisorSuppressionResource) ModelObject() interface{} {
	return &AdvisorSuppressionResourceModel{}
}

func (AdvisorSuppressionResource) ResourceType() string {
	return "azurerm_advisor_suppression"
}

func (r AdvisorSuppressionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Advisor.SuppressionsClient

			var model AdvisorSuppressionResourceModel

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := suppressions.NewScopedSuppressionID(model.ResourceID, model.RecommendationID, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			param := suppressions.SuppressionContract{
				Name: pointer.To(model.Name),
				Properties: &suppressions.SuppressionProperties{
					SuppressionId: pointer.To(model.SuppressionID),
				},
			}

			if model.TTL != "" {
				param.Properties.Ttl = pointer.To(model.TTL)
			}

			if _, err := client.Create(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (AdvisorSuppressionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Advisor.SuppressionsClient

			state := AdvisorSuppressionResourceModel{}

			id, err := suppressions.ParseScopedSuppressionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state.Name = id.SuppressionName
			state.ResourceID = id.ResourceUri
			state.RecommendationID = id.RecommendationId

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.TTL = pointer.From(props.Ttl)
					state.SuppressionID = pointer.From(props.SuppressionId)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (AdvisorSuppressionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Advisor.SuppressionsClient

			id, err := suppressions.ParseScopedSuppressionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (AdvisorSuppressionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return suppressions.ValidateScopedSuppressionID
}
