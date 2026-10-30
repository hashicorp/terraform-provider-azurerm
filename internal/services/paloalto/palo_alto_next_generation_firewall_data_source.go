// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package paloalto

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	firewalls "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/firewallresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (NextGenerationFirewallVHubLocalRulestackDataSource) ResourceType() string {
	return "azurerm_palo_alto_next_generation_firewall_virtual_hub_local_rulestack"
}

func (NextGenerationFirewallVHubLocalRulestackDataSource) ModelObject() interface{} {
	return &NextGenerationFirewallVHubLocalRuleStackModel{}
}

func (NextGenerationFirewallVHubLocalRulestackDataSource) Arguments() map[string]*pluginsdk.Schema {
	return nextGenerationFirewallDataSourceArguments()
}

func (NextGenerationFirewallVHubLocalRulestackDataSource) Attributes() map[string]*pluginsdk.Schema {
	attributes := nextGenerationFirewallDataSourceAttributes(schema.VHubNetworkProfileDataSourceSchema())
	attributes["rulestack_id"] = computedStringSchema()
	return attributes
}

func (NextGenerationFirewallVHubLocalRulestackDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var state NextGenerationFirewallVHubLocalRuleStackModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			id, model, err := retrieveNextGenerationFirewall(ctx, metadata, state.ResourceGroupName, state.Name)
			if err != nil {
				return err
			}

			if model != nil {
				props := model.Properties
				state.DNSSettings = schema.FlattenDNSSettings(props.DnsSettings)
				state.FrontEnd = schema.FlattenDestinationNAT(props.FrontEndSettings)
				state.MarketplaceOfferId = props.MarketplaceDetails.OfferId
				state.PlanId = props.PlanData.PlanId
				state.Tags = tags.Flatten(model.Tags)
				if props.AssociatedRulestack != nil {
					state.RuleStackId = pointer.From(props.AssociatedRulestack.ResourceId)
				}

				networkProfile, err := schema.FlattenNetworkProfileVHub(props.NetworkProfile)
				if err != nil {
					return fmt.Errorf("flattening Network Profile for %s: %+v", id, err)
				}
				state.NetworkProfile = []schema.NetworkProfileVHub{*networkProfile}
			}

			return metadata.Encode(&state)
		},
	}
}

func (NextGenerationFirewallVNetLocalRulestackDataSource) ResourceType() string {
	return "azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack"
}

func (NextGenerationFirewallVNetLocalRulestackDataSource) ModelObject() interface{} {
	return &NextGenerationFirewallVnetLocalRulestackModel{}
}

func (NextGenerationFirewallVNetLocalRulestackDataSource) Arguments() map[string]*pluginsdk.Schema {
	return nextGenerationFirewallDataSourceArguments()
}

func (NextGenerationFirewallVNetLocalRulestackDataSource) Attributes() map[string]*pluginsdk.Schema {
	attributes := nextGenerationFirewallDataSourceAttributes(schema.VnetNetworkProfileDataSourceSchema())
	attributes["rulestack_id"] = computedStringSchema()
	return attributes
}

func (NextGenerationFirewallVNetLocalRulestackDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var state NextGenerationFirewallVnetLocalRulestackModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			_, model, err := retrieveNextGenerationFirewall(ctx, metadata, state.ResourceGroupName, state.Name)
			if err != nil {
				return err
			}

			if model != nil {
				props := model.Properties
				state.DNSSettings = schema.FlattenDNSSettings(props.DnsSettings)
				state.FrontEnd = schema.FlattenDestinationNAT(props.FrontEndSettings)
				state.MarketplaceOfferId = props.MarketplaceDetails.OfferId
				state.NetworkProfile = schema.FlattenNetworkProfileVnet(props.NetworkProfile)
				state.PlanId = props.PlanData.PlanId
				state.Tags = tags.Flatten(model.Tags)
				if props.AssociatedRulestack != nil {
					state.RuleStackId = pointer.From(props.AssociatedRulestack.ResourceId)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (NextGenerationFirewallVHubPanoramaDataSource) ResourceType() string {
	return "azurerm_palo_alto_next_generation_firewall_virtual_hub_panorama"
}

func (NextGenerationFirewallVHubPanoramaDataSource) ModelObject() interface{} {
	return &NextGenerationFirewallVHubPanoramaResourceModel{}
}

func (NextGenerationFirewallVHubPanoramaDataSource) Arguments() map[string]*pluginsdk.Schema {
	return nextGenerationFirewallDataSourceArguments()
}

func (NextGenerationFirewallVHubPanoramaDataSource) Attributes() map[string]*pluginsdk.Schema {
	attributes := nextGenerationFirewallDataSourceAttributes(schema.VHubNetworkProfileDataSourceSchema())
	attributes["location"] = commonschema.LocationComputed()
	attributes["panorama"] = schema.PanoramaSchema()
	attributes["panorama_base64_config"] = computedStringSchema()
	return attributes
}

func (NextGenerationFirewallVHubPanoramaDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var state NextGenerationFirewallVHubPanoramaResourceModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			id, model, err := retrieveNextGenerationFirewall(ctx, metadata, state.ResourceGroupName, state.Name)
			if err != nil {
				return err
			}

			if model != nil {
				props := model.Properties
				state.Location = location.Normalize(model.Location)
				state.DNSSettings = schema.FlattenDNSSettings(props.DnsSettings)
				state.FrontEnd = schema.FlattenDestinationNAT(props.FrontEndSettings)
				state.MarketplaceOfferId = props.MarketplaceDetails.OfferId
				state.PlanId = props.PlanData.PlanId
				state.Tags = tags.Flatten(model.Tags)

				networkProfile, err := schema.FlattenNetworkProfileVHub(props.NetworkProfile)
				if err != nil {
					return fmt.Errorf("flattening Network Profile for %s: %+v", id, err)
				}
				state.NetworkProfile = []schema.NetworkProfileVHub{*networkProfile}

				flattenPanoramaConfig(props.PanoramaConfig, &state.PanoramaBase64Config, &state.PanoramaConfig)
			}

			return metadata.Encode(&state)
		},
	}
}

func (NextGenerationFirewallVNetPanoramaDataSource) ResourceType() string {
	return "azurerm_palo_alto_next_generation_firewall_virtual_network_panorama"
}

func (NextGenerationFirewallVNetPanoramaDataSource) ModelObject() interface{} {
	return &NextGenerationFirewallVnetPanoramaModel{}
}

func (NextGenerationFirewallVNetPanoramaDataSource) Arguments() map[string]*pluginsdk.Schema {
	return nextGenerationFirewallDataSourceArguments()
}

func (NextGenerationFirewallVNetPanoramaDataSource) Attributes() map[string]*pluginsdk.Schema {
	attributes := nextGenerationFirewallDataSourceAttributes(schema.VnetNetworkProfileDataSourceSchema())
	attributes["location"] = commonschema.LocationComputed()
	attributes["panorama"] = schema.PanoramaSchema()
	attributes["panorama_base64_config"] = computedStringSchema()
	return attributes
}

func (NextGenerationFirewallVNetPanoramaDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var state NextGenerationFirewallVnetPanoramaModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			_, model, err := retrieveNextGenerationFirewall(ctx, metadata, state.ResourceGroupName, state.Name)
			if err != nil {
				return err
			}

			if model != nil {
				props := model.Properties
				state.Location = location.Normalize(model.Location)
				state.DNSSettings = schema.FlattenDNSSettings(props.DnsSettings)
				state.FrontEnd = schema.FlattenDestinationNAT(props.FrontEndSettings)
				state.MarketplaceOfferId = props.MarketplaceDetails.OfferId
				state.NetworkProfile = schema.FlattenNetworkProfileVnet(props.NetworkProfile)
				state.PlanId = props.PlanData.PlanId
				state.Tags = tags.Flatten(model.Tags)
				flattenPanoramaConfig(props.PanoramaConfig, &state.PanoramaBase64Config, &state.PanoramaConfig)
			}

			return metadata.Encode(&state)
		},
	}
}

func (NextGenerationFirewallVHubStrataCloudManagerDataSource) ResourceType() string {
	return "azurerm_palo_alto_next_generation_firewall_virtual_hub_strata_cloud_manager"
}

func (NextGenerationFirewallVHubStrataCloudManagerDataSource) ModelObject() interface{} {
	return &NextGenerationFirewallVHubStrataCloudManagerModel{}
}

func (NextGenerationFirewallVHubStrataCloudManagerDataSource) Arguments() map[string]*pluginsdk.Schema {
	return nextGenerationFirewallDataSourceArguments()
}

func (NextGenerationFirewallVHubStrataCloudManagerDataSource) Attributes() map[string]*pluginsdk.Schema {
	attributes := nextGenerationFirewallDataSourceAttributes(schema.VHubNetworkProfileDataSourceSchema())
	attributes["identity"] = userAssignedIdentityDataSourceSchema()
	attributes["location"] = commonschema.LocationComputed()
	attributes["strata_cloud_manager_tenant_name"] = computedStringSchema()
	return attributes
}

func (NextGenerationFirewallVHubStrataCloudManagerDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var state NextGenerationFirewallVHubStrataCloudManagerModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			id, model, err := retrieveNextGenerationFirewall(ctx, metadata, state.ResourceGroupName, state.Name)
			if err != nil {
				return err
			}

			if model != nil {
				props := model.Properties
				state.Location = location.Normalize(model.Location)
				state.DNSSettings = schema.FlattenDNSSettings(props.DnsSettings)
				state.FrontEnd = schema.FlattenDestinationNAT(props.FrontEndSettings)
				state.MarketplaceOfferId = props.MarketplaceDetails.OfferId
				state.PlanId = props.PlanData.PlanId
				state.Tags = tags.Flatten(model.Tags)

				networkProfile, err := schema.FlattenNetworkProfileVHub(props.NetworkProfile)
				if err != nil {
					return fmt.Errorf("flattening Network Profile for %s: %+v", id, err)
				}
				state.NetworkProfile = []schema.NetworkProfileVHub{*networkProfile}

				if err := flattenStrataCloudManagerConfig(model, &state.StrataCloudManagerTenantName, &state.Identity); err != nil {
					return err
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (NextGenerationFirewallVNetStrataCloudManagerDataSource) ResourceType() string {
	return "azurerm_palo_alto_next_generation_firewall_virtual_network_strata_cloud_manager"
}

func (NextGenerationFirewallVNetStrataCloudManagerDataSource) ModelObject() interface{} {
	return &NextGenerationFirewallVNetStrataCloudManagerModel{}
}

func (NextGenerationFirewallVNetStrataCloudManagerDataSource) Arguments() map[string]*pluginsdk.Schema {
	return nextGenerationFirewallDataSourceArguments()
}

func (NextGenerationFirewallVNetStrataCloudManagerDataSource) Attributes() map[string]*pluginsdk.Schema {
	attributes := nextGenerationFirewallDataSourceAttributes(schema.VnetNetworkProfileDataSourceSchema())
	attributes["identity"] = userAssignedIdentityDataSourceSchema()
	attributes["location"] = commonschema.LocationComputed()
	attributes["strata_cloud_manager_tenant_name"] = computedStringSchema()
	return attributes
}

func (NextGenerationFirewallVNetStrataCloudManagerDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var state NextGenerationFirewallVNetStrataCloudManagerModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			_, model, err := retrieveNextGenerationFirewall(ctx, metadata, state.ResourceGroupName, state.Name)
			if err != nil {
				return err
			}

			if model != nil {
				props := model.Properties
				state.Location = location.Normalize(model.Location)
				state.DNSSettings = schema.FlattenDNSSettings(props.DnsSettings)
				state.FrontEnd = schema.FlattenDestinationNAT(props.FrontEndSettings)
				state.MarketplaceOfferId = props.MarketplaceDetails.OfferId
				state.NetworkProfile = schema.FlattenNetworkProfileVnet(props.NetworkProfile)
				state.PlanId = props.PlanData.PlanId
				state.Tags = tags.Flatten(model.Tags)

				if err := flattenStrataCloudManagerConfig(model, &state.StrataCloudManagerTenantName, &state.Identity); err != nil {
					return err
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func nextGenerationFirewallDataSourceArguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.NextGenerationFirewallName,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func nextGenerationFirewallDataSourceAttributes(networkProfile *pluginsdk.Schema) map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"destination_nat":      schema.DestinationNATDataSourceSchema(),
		"dns_settings":         schema.DNSSettingsDataSourceSchema(),
		"marketplace_offer_id": computedStringSchema(),
		"network_profile":      networkProfile,
		"plan_id":              computedStringSchema(),
		"tags":                 commonschema.TagsDataSource(),
	}
}

func computedStringSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeString,
		Computed: true,
	}
}

func userAssignedIdentityDataSourceSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"identity_ids": {
					Type:     pluginsdk.TypeSet,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},
				"type": computedStringSchema(),
			},
		},
	}
}

func retrieveNextGenerationFirewall(ctx context.Context, metadata sdk.ResourceMetaData, resourceGroupName, name string) (firewalls.FirewallId, *firewalls.FirewallResource, error) {
	client := metadata.Client.PaloAlto.FirewallResources
	id := firewalls.NewFirewallID(metadata.Client.Account.SubscriptionId, resourceGroupName, name)

	result, err := client.FirewallsGet(ctx, id)
	if err != nil {
		if response.WasNotFound(result.HttpResponse) {
			return id, nil, fmt.Errorf("%s was not found", id)
		}
		return id, nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	metadata.SetID(id)
	return id, result.Model, nil
}

func flattenPanoramaConfig(input *firewalls.PanoramaConfig, base64Config *string, panoramaConfig *[]schema.Panorama) {
	if input == nil {
		return
	}

	*base64Config = input.ConfigString
	*panoramaConfig = []schema.Panorama{{
		Name:            pointer.From(input.CgName),
		DeviceGroupName: pointer.From(input.DgName),
		HostName:        pointer.From(input.HostName),
		PanoramaServer:  pointer.From(input.PanoramaServer),
		PanoramaServer2: pointer.From(input.PanoramaServer2),
		TplName:         pointer.From(input.TplName),
		VMAuthKey:       pointer.From(input.VMAuthKey),
	}}
}

func flattenStrataCloudManagerConfig(input *firewalls.FirewallResource, tenantName *string, outputIdentity *[]identity.ModelUserAssigned) error {
	if config := input.Properties.StrataCloudManagerConfig; config != nil {
		*tenantName = config.CloudManagerName
	}

	flattenedIdentity, err := identity.FlattenUserAssignedMapToModel(input.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	*outputIdentity = pointer.From(flattenedIdentity)
	return nil
}
