package containers

// NOTE: this file is generated - manual changes will be overwritten.
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.
import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-04-01/fleets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.Resource           = KubernetesFleetManagerResource{}
	_ sdk.ResourceWithUpdate = KubernetesFleetManagerResource{}
)

type KubernetesFleetManagerResource struct{}

func (r KubernetesFleetManagerResource) ModelObject() interface{} {
	return &KubernetesFleetManagerResourceSchema{}
}

type KubernetesFleetManagerResourceSchema struct {
	Location          string                   `tfschema:"location"`
	Name              string                   `tfschema:"name"`
	ResourceGroupName string                   `tfschema:"resource_group_name"`
	HubProfile        []FleetManagerHubProfile `tfschema:"hub_profile"`
	Tags              map[string]interface{}   `tfschema:"tags"`
}

type FleetManagerHubProfile struct {
	DnsPrefix         string `tfschema:"dns_prefix"`
	Fqdn              string `tfschema:"fqdn"`
	KubernetesVersion string `tfschema:"kubernetes_version"`
	PortalFqdn        string `tfschema:"portal_fqdn"`
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
						ForceNew: true,
						Type:     pluginsdk.TypeString,
						ValidateFunc: validation.All(
							validation.StringLenBetween(1, 54),
							validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9]$|^[a-zA-Z0-9][a-zA-Z0-9-]{0,52}[a-zA-Z0-9]$`), "must match the pattern ^[a-zA-Z0-9]$|^[a-zA-Z0-9][a-zA-Z0-9-]{0,52}[a-zA-Z0-9]$"),
						),
					},
					"fqdn": {
						Computed: true,
						Type:     pluginsdk.TypeString,
					},
					"kubernetes_version": {
						Computed: true,
						Type:     pluginsdk.TypeString,
					},
					"portal_fqdn": {
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
			client := metadata.Client.ContainerService.V20240401.Fleets

			var config KubernetesFleetManagerResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId

			id := fleets.NewFleetID(subscriptionId, config.ResourceGroupName, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			var payload fleets.Fleet
			r.mapKubernetesFleetManagerResourceSchemaToFleet(config, &payload)

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, fleets.DefaultCreateOrUpdateOperationOptions(), metadata.SetIDCallback(&id)); err != nil {
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
			client := metadata.Client.ContainerService.V20240401.Fleets
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
				r.mapFleetToKubernetesFleetManagerResourceSchema(*model, &schema)
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r KubernetesFleetManagerResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerService.V20240401.Fleets

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
			client := metadata.Client.ContainerService.V20240401.Fleets

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

			r.mapKubernetesFleetManagerResourceSchemaToFleet(config, &payload)

			if err := client.CreateOrUpdateThenPoll(ctx, *id, payload, fleets.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r KubernetesFleetManagerResource) mapKubernetesFleetManagerResourceSchemaToFleet(input KubernetesFleetManagerResourceSchema, output *fleets.Fleet) {
	output.Location = location.Normalize(input.Location)
	output.Tags = tags.Expand(input.Tags)

	if output.Properties == nil {
		output.Properties = &fleets.FleetProperties{}
	}

	if len(input.HubProfile) > 0 {
		output.Properties.HubProfile = expandFleetManagerHubProfile(input.HubProfile)
	}
}

func (r KubernetesFleetManagerResource) mapFleetToKubernetesFleetManagerResourceSchema(input fleets.Fleet, output *KubernetesFleetManagerResourceSchema) {
	output.Location = location.Normalize(input.Location)
	output.Tags = tags.Flatten(input.Tags)

	if input.Properties == nil {
		input.Properties = &fleets.FleetProperties{}
	}

	output.HubProfile = flattenFleetManagerHubProfile(input.Properties.HubProfile)
}

func expandFleetManagerHubProfile(input []FleetManagerHubProfile) *fleets.FleetHubProfile {
	if len(input) == 0 {
		return nil
	}

	return &fleets.FleetHubProfile{
		DnsPrefix: pointer.To(input[0].DnsPrefix),
	}
}

func flattenFleetManagerHubProfile(input *fleets.FleetHubProfile) []FleetManagerHubProfile {
	if input == nil {
		return []FleetManagerHubProfile{}
	}

	return []FleetManagerHubProfile{
		{
			DnsPrefix:         pointer.From(input.DnsPrefix),
			Fqdn:              pointer.From(input.Fqdn),
			KubernetesVersion: pointer.From(input.KubernetesVersion),
			PortalFqdn:        pointer.From(input.PortalFqdn),
		},
	}
}
