// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2019-09-01/querypacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LogAnalyticsQueryPackModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	Tags              map[string]string `tfschema:"tags"`
}

type LogAnalyticsQueryPackResource struct{}

var _ sdk.ResourceWithUpdate = LogAnalyticsQueryPackResource{}

func (r LogAnalyticsQueryPackResource) ResourceType() string {
	return "azurerm_log_analytics_query_pack"
}

func (r LogAnalyticsQueryPackResource) ModelObject() interface{} {
	return &LogAnalyticsQueryPackModel{}
}

func (r LogAnalyticsQueryPackResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return querypacks.ValidateQueryPackID
}

func (r LogAnalyticsQueryPackResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"tags": commonschema.Tags(),
	}
}

func (r LogAnalyticsQueryPackResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r LogAnalyticsQueryPackResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model LogAnalyticsQueryPackModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.LogAnalytics.QueryPacksClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := querypacks.NewQueryPackID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.QueryPacksGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &querypacks.LogAnalyticsQueryPack{
				Location:   location.Normalize(model.Location),
				Properties: querypacks.LogAnalyticsQueryPackProperties{},
				Tags:       &model.Tags,
			}

			if resp, err := client.QueryPacksCreateOrUpdate(ctx, id, *properties); err != nil {
				// update check logic once the issue https://github.com/Azure/azure-rest-api-specs/issues/19603 is fixed
				if !response.WasStatusCode(resp.HttpResponse, http.StatusCreated) {
					return fmt.Errorf("creating %s: %+v", id, err)
				}
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r LogAnalyticsQueryPackResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LogAnalytics.QueryPacksClient

			id, err := querypacks.ParseQueryPackID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model LogAnalyticsQueryPackModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.QueryPacksGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
			}

			if resp, err := client.QueryPacksCreateOrUpdate(ctx, *id, *properties); err != nil {
				// update check logic once the issue https://github.com/Azure/azure-rest-api-specs/issues/19603 is fixed
				if !response.WasStatusCode(resp.HttpResponse, http.StatusCreated) {
					return fmt.Errorf("updating %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}

func (r LogAnalyticsQueryPackResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LogAnalytics.QueryPacksClient

			id, err := querypacks.ParseQueryPackID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.QueryPacksGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := LogAnalyticsQueryPackModel{
				Name:              id.QueryPackName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.NormalizeNilable(utils.String(model.Location)),
			}

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r LogAnalyticsQueryPackResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LogAnalytics.QueryPacksClient

			id, err := querypacks.ParseQueryPackID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.QueryPacksDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
