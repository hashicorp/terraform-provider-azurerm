// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package paloalto

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	firewalls "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/firewallresources"
	metricsobjectfirewall "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/metricsobjectfirewallresources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NextGenerationFirewallMetricsResource struct{}

type NextGenerationFirewallMetricsModel struct {
	FirewallID                          string `tfschema:"firewall_id"`
	ApplicationInsightsConnectionString string `tfschema:"application_insights_connection_string"`
	ApplicationInsightsResourceID       string `tfschema:"application_insights_resource_id"`
	PanEtag                             string `tfschema:"pan_etag"`
}

var _ sdk.ResourceWithUpdate = NextGenerationFirewallMetricsResource{}

func (r NextGenerationFirewallMetricsResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return firewalls.ValidateFirewallID
}

func (r NextGenerationFirewallMetricsResource) ResourceType() string {
	return "azurerm_palo_alto_next_generation_firewall_metrics"
}

func (r NextGenerationFirewallMetricsResource) ModelObject() interface{} {
	return &NextGenerationFirewallMetricsModel{}
}

func (r NextGenerationFirewallMetricsResource) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"firewall_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: firewalls.ValidateFirewallID,
		},

		"application_insights_connection_string": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"application_insights_resource_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: azure.ValidateResourceID,
		},
	}
}

func (r NextGenerationFirewallMetricsResource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"pan_etag": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r NextGenerationFirewallMetricsResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.MetricsObjectFirewallResources

			var model NextGenerationFirewallMetricsModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			firewallId, err := firewalls.ParseFirewallID(model.FirewallID)
			if err != nil {
				return err
			}

			metricsFirewallId := metricsobjectfirewall.NewFirewallID(firewallId.SubscriptionId, firewallId.ResourceGroupName, firewallId.FirewallName)

			existing, err := client.MetricsObjectFirewallGet(ctx, metricsFirewallId)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing Metrics for %s: %+v", metricsFirewallId, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), firewallId)
			}

			input := metricsobjectfirewall.MetricsObjectFirewallResource{
				Properties: metricsobjectfirewall.MetricsObject{
					ApplicationInsightsConnectionString: model.ApplicationInsightsConnectionString,
					ApplicationInsightsResourceId:       model.ApplicationInsightsResourceID,
				},
			}

			if err = client.MetricsObjectFirewallCreateOrUpdateThenPoll(ctx, metricsFirewallId, input); err != nil {
				return fmt.Errorf("creating Metrics for %s: %+v", metricsFirewallId, err)
			}

			metadata.SetID(firewallId)

			return nil
		},
	}
}

func (r NextGenerationFirewallMetricsResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.MetricsObjectFirewallResources

			firewallId, err := firewalls.ParseFirewallID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metricsFirewallId := metricsobjectfirewall.NewFirewallID(firewallId.SubscriptionId, firewallId.ResourceGroupName, firewallId.FirewallName)

			existing, err := client.MetricsObjectFirewallGet(ctx, metricsFirewallId)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(firewallId)
				}
				return fmt.Errorf("reading Metrics for %s: %+v", metricsFirewallId, err)
			}

			state := NextGenerationFirewallMetricsModel{
				FirewallID: firewallId.ID(),
			}

			if model := existing.Model; model != nil {
				props := model.Properties
				// Azure may redact the connection string on GET; preserve the config value in state
				if props.ApplicationInsightsConnectionString != "" {
					state.ApplicationInsightsConnectionString = props.ApplicationInsightsConnectionString
				}
				state.ApplicationInsightsResourceID = props.ApplicationInsightsResourceId
				state.PanEtag = pointer.From(props.PanEtag)
			}

			return metadata.Encode(&state)
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
				return err
			}

			firewallId, err := firewalls.ParseFirewallID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metricsFirewallId := metricsobjectfirewall.NewFirewallID(firewallId.SubscriptionId, firewallId.ResourceGroupName, firewallId.FirewallName)

			existing, err := client.MetricsObjectFirewallGet(ctx, metricsFirewallId)
			if err != nil {
				return fmt.Errorf("retrieving Metrics for %s: %+v", metricsFirewallId, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving Metrics for %s: `model` was nil", metricsFirewallId)
			}

			update := *existing.Model

			if metadata.ResourceData.HasChange("application_insights_connection_string") {
				update.Properties.ApplicationInsightsConnectionString = model.ApplicationInsightsConnectionString
			}

			if metadata.ResourceData.HasChange("application_insights_resource_id") {
				update.Properties.ApplicationInsightsResourceId = model.ApplicationInsightsResourceID
			}

			if err = client.MetricsObjectFirewallCreateOrUpdateThenPoll(ctx, metricsFirewallId, update); err != nil {
				return fmt.Errorf("updating Metrics for %s: %+v", metricsFirewallId, err)
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

			firewallId, err := firewalls.ParseFirewallID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metricsFirewallId := metricsobjectfirewall.NewFirewallID(firewallId.SubscriptionId, firewallId.ResourceGroupName, firewallId.FirewallName)

			if err = client.MetricsObjectFirewallDeleteThenPoll(ctx, metricsFirewallId); err != nil {
				return fmt.Errorf("deleting Metrics for %s: %+v", metricsFirewallId, err)
			}

			return nil
		},
	}
}
