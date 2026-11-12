// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package paloalto

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/metricsobjectfirewallresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -parent-id firewall_id -test-expect-non-empty

type NextGenerationFirewallMetricsResource struct{}

type NextGenerationFirewallMetricsModel struct {
	FirewallID                          string `tfschema:"firewall_id"`
	ApplicationInsightsConnectionString string `tfschema:"application_insights_connection_string"`
	ApplicationInsightsID               string `tfschema:"application_insights_id"`
}

var (
	_ sdk.ResourceWithUpdate               = NextGenerationFirewallMetricsResource{}
	_ sdk.ResourceWithIdentityTypeOverride = NextGenerationFirewallMetricsResource{}
)

func (r NextGenerationFirewallMetricsResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return metricsobjectfirewallresources.ValidateFirewallID
}

func (r NextGenerationFirewallMetricsResource) Identity() resourceids.ResourceId {
	return &metricsobjectfirewallresources.FirewallId{}
}

func (r NextGenerationFirewallMetricsResource) IdentityType() pluginsdk.ResourceTypeForIdentity {
	return pluginsdk.ResourceTypeForIdentityVirtual
}

func (r NextGenerationFirewallMetricsResource) ResourceType() string {
	return "azurerm_palo_alto_next_generation_firewall_metrics"
}

func (r NextGenerationFirewallMetricsResource) ModelObject() any {
	return &NextGenerationFirewallMetricsModel{}
}

func (r NextGenerationFirewallMetricsResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"firewall_id": commonschema.ResourceIDReferenceRequiredForceNew(&metricsobjectfirewallresources.FirewallId{}),

		"application_insights_connection_string": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"application_insights_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: componentsapis.ValidateComponentID,
		},
	}
}

func (r NextGenerationFirewallMetricsResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r NextGenerationFirewallMetricsResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.MetricsObjectFirewallResources

			var model NextGenerationFirewallMetricsModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := metricsobjectfirewallresources.ParseFirewallID(model.FirewallID)
			if err != nil {
				return err
			}

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.MetricsObjectFirewallGet(ctx, *id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing Metrics Object for %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			input := metricsobjectfirewallresources.MetricsObjectFirewallResource{
				Properties: metricsobjectfirewallresources.MetricsObject{
					ApplicationInsightsConnectionString: model.ApplicationInsightsConnectionString,
					ApplicationInsightsResourceId:       model.ApplicationInsightsID,
				},
			}

			if err := client.MetricsObjectFirewallCreateOrUpdateCallbackThenPoll(ctx, *id, input, metadata.SetIDAndIdentityWithTypeCallback(id, pluginsdk.ResourceTypeForIdentityVirtual)); err != nil {
				return fmt.Errorf("creating Metrics Object for %s: %+v", id, err)
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, id, pluginsdk.ResourceTypeForIdentityVirtual)
		},
	}
}

func (r NextGenerationFirewallMetricsResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.MetricsObjectFirewallResources

			id, err := metricsobjectfirewallresources.ParseFirewallID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.MetricsObjectFirewallGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving Metrics Object for %s: %+v", id, err)
			}

			return r.flatten(metadata, id, existing.Model)
		},
	}
}

func (r NextGenerationFirewallMetricsResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.MetricsObjectFirewallResources

			var model NextGenerationFirewallMetricsModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := metricsobjectfirewallresources.ParseFirewallID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.MetricsObjectFirewallGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving Metrics Object for %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving Metrics Object for %s: `model` was nil", id)
			}

			payload := *existing.Model

			// the API may redact the connection string on GET, so it is always sent from config to keep the PUT payload valid
			payload.Properties.ApplicationInsightsConnectionString = model.ApplicationInsightsConnectionString

			if metadata.ResourceData.HasChange("application_insights_id") {
				payload.Properties.ApplicationInsightsResourceId = model.ApplicationInsightsID
			}

			if err := client.MetricsObjectFirewallCreateOrUpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating Metrics Object for %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r NextGenerationFirewallMetricsResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.MetricsObjectFirewallResources

			id, err := metricsobjectfirewallresources.ParseFirewallID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.MetricsObjectFirewallDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting Metrics Object for %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r NextGenerationFirewallMetricsResource) flatten(metadata sdk.ResourceMetaData, id *metricsobjectfirewallresources.FirewallId, model *metricsobjectfirewallresources.MetricsObjectFirewallResource) error {
	state := NextGenerationFirewallMetricsModel{
		FirewallID: id.ID(),
	}

	if model != nil {
		props := model.Properties

		// the API may redact the connection string on GET, so default to the value currently in state
		state.ApplicationInsightsConnectionString = metadata.ResourceData.Get("application_insights_connection_string").(string)
		if props.ApplicationInsightsConnectionString != "" {
			state.ApplicationInsightsConnectionString = props.ApplicationInsightsConnectionString
		}

		if props.ApplicationInsightsResourceId != "" {
			applicationInsightsId, err := componentsapis.ParseComponentIDInsensitively(props.ApplicationInsightsResourceId)
			if err != nil {
				return err
			}
			state.ApplicationInsightsID = applicationInsightsId.ID()
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id, pluginsdk.ResourceTypeForIdentityVirtual); err != nil {
		return err
	}

	return metadata.Encode(&state)
}
