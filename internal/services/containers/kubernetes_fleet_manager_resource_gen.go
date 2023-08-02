package containers

// NOTE: this file is generated - manual changes will be overwritten.
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.
import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2022-09-02-preview/fleets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = KubernetesFleetManagerResource{}
var _ sdk.ResourceWithUpdate = KubernetesFleetManagerResource{}

type KubernetesFleetManagerResource struct{}

func (r KubernetesFleetManagerResource) ModelObject() interface{} {
	return &KubernetesFleetManagerResourceSchema{}
}

type KubernetesFleetManagerResourceSchema struct {
	HubProfile        []KubernetesFleetManagerResourceFleetHubProfileSchema `tfschema:"hub_profile"`
	Location          string                                                `tfschema:"location"`
	Name              string                                                `tfschema:"name"`
	ResourceGroupName string                                                `tfschema:"resource_group_name"`
	Tags              map[string]interface{}                                `tfschema:"tags"`
}

func (r KubernetesFleetManagerResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return fleets.ValidateFleetID
}
func (r KubernetesFleetManagerResource) ResourceType() string {
	return "azurerm_kubernetes_fleet_manager"
}
func (r KubernetesFleetManagerResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),
		"name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"resource_group_name": commonschema.ResourceGroupName(),
		"hub_profile": {
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"dns_prefix": {
						Required: true,
						Type:     pluginsdk.TypeString,
					},
					"fqdn": {
						Computed: true,
						Type:     pluginsdk.TypeString,
					},
					"kubernetes_version": {
						Computed: true,
						Type:     pluginsdk.TypeString,
					},
				},
			},
			ForceNew: true,
			MaxItems: 1,
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
		"tags": commonschema.Tags(),
	}
}
func (r KubernetesFleetManagerResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
func (r KubernetesFleetManagerResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerService.V20220902Preview.Fleets

			var config KubernetesFleetManagerResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId

			id := fleets.NewFleetID(subscriptionId, config.ResourceGroupName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			var payload fleets.Fleet
			if err := r.mapKubernetesFleetManagerResourceSchemaToFleet(config, &payload); err != nil {
				return fmt.Errorf("mapping schema model to sdk model: %+v", err)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload, fleets.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
func (r KubernetesFleetManagerResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerService.V20220902Preview.Fleets
			schema := KubernetesFleetManagerResourceSchema{}

			id, err := fleets.ParseFleetID(metadata.ResourceData.Id())
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

			if model := resp.Model; model != nil {
				schema.Name = id.FleetName
				schema.ResourceGroupName = id.ResourceGroupName
				if err := r.mapFleetToKubernetesFleetManagerResourceSchema(*model, &schema); err != nil {
					return fmt.Errorf("flattening model: %+v", err)
				}
			}

			return metadata.Encode(&schema)
		},
	}
}
func (r KubernetesFleetManagerResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerService.V20220902Preview.Fleets

			id, err := fleets.ParseFleetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, fleets.DefaultDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
func (r KubernetesFleetManagerResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerService.V20220902Preview.Fleets

			id, err := fleets.ParseFleetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config KubernetesFleetManagerResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving existing %s: %+v", *id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving existing %s: properties was nil", *id)
			}
			payload := *existing.Model

			if err := r.mapKubernetesFleetManagerResourceSchemaToFleet(config, &payload); err != nil {
				return fmt.Errorf("mapping schema model to sdk model: %+v", err)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, payload, fleets.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

type KubernetesFleetManagerResourceFleetHubProfileSchema struct {
	DnsPrefix         string `tfschema:"dns_prefix"`
	Fqdn              string `tfschema:"fqdn"`
	KubernetesVersion string `tfschema:"kubernetes_version"`
}

func (r KubernetesFleetManagerResource) mapKubernetesFleetManagerResourceFleetHubProfileSchemaToFleetHubProfile(input KubernetesFleetManagerResourceFleetHubProfileSchema, output *fleets.FleetHubProfile) error {
	output.DnsPrefix = input.DnsPrefix

	return nil
}

func (r KubernetesFleetManagerResource) mapFleetHubProfileToKubernetesFleetManagerResourceFleetHubProfileSchema(input fleets.FleetHubProfile, output *KubernetesFleetManagerResourceFleetHubProfileSchema) error {
	output.DnsPrefix = input.DnsPrefix
	output.Fqdn = pointer.From(input.Fqdn)
	output.KubernetesVersion = pointer.From(input.KubernetesVersion)
	return nil
}

func (r KubernetesFleetManagerResource) mapKubernetesFleetManagerResourceSchemaToFleetProperties(input KubernetesFleetManagerResourceSchema, output *fleets.FleetProperties) error {
	if len(input.HubProfile) > 0 {
		if err := r.mapKubernetesFleetManagerResourceFleetHubProfileSchemaToFleetProperties(input.HubProfile[0], output); err != nil {
			return err
		}
	}
	return nil
}

func (r KubernetesFleetManagerResource) mapFleetPropertiesToKubernetesFleetManagerResourceSchema(input fleets.FleetProperties, output *KubernetesFleetManagerResourceSchema) error {
	tmpHubProfile := &KubernetesFleetManagerResourceFleetHubProfileSchema{}
	if err := r.mapFleetPropertiesToKubernetesFleetManagerResourceFleetHubProfileSchema(input, tmpHubProfile); err != nil {
		return err
	} else {
		output.HubProfile = make([]KubernetesFleetManagerResourceFleetHubProfileSchema, 0)
		output.HubProfile = append(output.HubProfile, *tmpHubProfile)
	}
	return nil
}

func (r KubernetesFleetManagerResource) mapKubernetesFleetManagerResourceSchemaToFleet(input KubernetesFleetManagerResourceSchema, output *fleets.Fleet) error {
	output.Location = location.Normalize(input.Location)
	output.Tags = tags.Expand(input.Tags)

	if output.Properties == nil {
		output.Properties = &fleets.FleetProperties{}
	}
	if err := r.mapKubernetesFleetManagerResourceSchemaToFleetProperties(input, output.Properties); err != nil {
		return fmt.Errorf("mapping Schema to SDK Field %q / Model %q: %+v", "FleetProperties", "Properties", err)
	}

	return nil
}

func (r KubernetesFleetManagerResource) mapFleetToKubernetesFleetManagerResourceSchema(input fleets.Fleet, output *KubernetesFleetManagerResourceSchema) error {
	output.Location = location.Normalize(input.Location)
	output.Tags = tags.Flatten(input.Tags)

	if input.Properties == nil {
		input.Properties = &fleets.FleetProperties{}
	}
	if err := r.mapFleetPropertiesToKubernetesFleetManagerResourceSchema(*input.Properties, output); err != nil {
		return fmt.Errorf("mapping SDK Field %q / Model %q to Schema: %+v", "FleetProperties", "Properties", err)
	}

	return nil
}

func (r KubernetesFleetManagerResource) mapKubernetesFleetManagerResourceFleetHubProfileSchemaToFleetProperties(input KubernetesFleetManagerResourceFleetHubProfileSchema, output *fleets.FleetProperties) error {

	if output.HubProfile == nil {
		output.HubProfile = &fleets.FleetHubProfile{}
	}
	if err := r.mapKubernetesFleetManagerResourceFleetHubProfileSchemaToFleetHubProfile(input, output.HubProfile); err != nil {
		return fmt.Errorf("mapping Schema to SDK Field %q / Model %q: %+v", "FleetHubProfile", "HubProfile", err)
	}

	return nil
}

func (r KubernetesFleetManagerResource) mapFleetPropertiesToKubernetesFleetManagerResourceFleetHubProfileSchema(input fleets.FleetProperties, output *KubernetesFleetManagerResourceFleetHubProfileSchema) error {

	if input.HubProfile == nil {
		input.HubProfile = &fleets.FleetHubProfile{}
	}
	if err := r.mapFleetHubProfileToKubernetesFleetManagerResourceFleetHubProfileSchema(*input.HubProfile, output); err != nil {
		return fmt.Errorf("mapping SDK Field %q / Model %q to Schema: %+v", "FleetHubProfile", "HubProfile", err)
	}

	return nil
}
