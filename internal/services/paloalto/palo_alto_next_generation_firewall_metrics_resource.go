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
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/firewallresources"
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
	return firewallresources.ValidateFirewallID
}

func (r NextGenerationFirewallMetricsResource) Identity() resourceids.ResourceId {
	return &firewallresources.FirewallId{}
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
		"firewall_id": commonschema.ResourceIDReferenceRequiredForceNew(&firewallresources.FirewallId{}),

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

			firewallId, err := firewallresources.ParseFirewallID(model.FirewallID)
			if err != nil {
				return err
			}

			metricsFirewallId := metricsobjectfirewallresources.NewFirewallID(firewallId.SubscriptionId, firewallId.ResourceGroupName, firewallId.FirewallName)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.MetricsObjectFirewallGet(ctx, metricsFirewallId)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", metricsFirewallId, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), firewallId)
				}
			}

			input := metricsobjectfirewallresources.MetricsObjectFirewallResource{
				Properties: metricsobjectfirewallresources.MetricsObject{
					ApplicationInsightsConnectionString: model.ApplicationInsightsConnectionString,
					ApplicationInsightsResourceId:       model.ApplicationInsightsID,
				},
			}

			if err := client.MetricsObjectFirewallCreateOrUpdateCallbackThenPoll(ctx, metricsFirewallId, input, metadata.SetIDAndIdentityWithTypeCallback(firewallId, pluginsdk.ResourceTypeForIdentityVirtual)); err != nil {
				return fmt.Errorf("creating %s: %+v", metricsFirewallId, err)
			}

			metadata.SetID(firewallId)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, firewallId, pluginsdk.ResourceTypeForIdentityVirtual)
		},
	}
}

func (r NextGenerationFirewallMetricsResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.MetricsObjectFirewallResources

			firewallId, err := firewallresources.ParseFirewallID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metricsFirewallId := metricsobjectfirewallresources.NewFirewallID(firewallId.SubscriptionId, firewallId.ResourceGroupName, firewallId.FirewallName)

			existing, err := client.MetricsObjectFirewallGet(ctx, metricsFirewallId)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(firewallId)
				}
				return fmt.Errorf("retrieving %s: %+v", metricsFirewallId, err)
			}

			return r.flatten(metadata, firewallId, existing.Model)
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

			firewallId, err := firewallresources.ParseFirewallID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metricsFirewallId := metricsobjectfirewallresources.NewFirewallID(firewallId.SubscriptionId, firewallId.ResourceGroupName, firewallId.FirewallName)

			existing, err := client.MetricsObjectFirewallGet(ctx, metricsFirewallId)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", metricsFirewallId, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", metricsFirewallId)
			}

			payload := *existing.Model

			// the API may redact the connection string on GET, so it is always sent from config to keep the PUT payload valid
			payload.Properties.ApplicationInsightsConnectionString = model.ApplicationInsightsConnectionString

			if metadata.ResourceData.HasChange("application_insights_id") {
				payload.Properties.ApplicationInsightsResourceId = model.ApplicationInsightsID
			}

			if err := client.MetricsObjectFirewallCreateOrUpdateThenPoll(ctx, metricsFirewallId, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", metricsFirewallId, err)
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

			firewallId, err := firewallresources.ParseFirewallID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metricsFirewallId := metricsobjectfirewallresources.NewFirewallID(firewallId.SubscriptionId, firewallId.ResourceGroupName, firewallId.FirewallName)

			if err := client.MetricsObjectFirewallDeleteThenPoll(ctx, metricsFirewallId); err != nil {
				return fmt.Errorf("deleting %s: %+v", metricsFirewallId, err)
			}

			return nil
		},
	}
}

func (r NextGenerationFirewallMetricsResource) flatten(metadata sdk.ResourceMetaData, id *firewallresources.FirewallId, model *metricsobjectfirewallresources.MetricsObjectFirewallResource) error {
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
