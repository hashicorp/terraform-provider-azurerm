// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-04-01/fleets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.DataSource = KubernetesFleetManagerDataSource{}

type KubernetesFleetManagerDataSource struct{}

type KubernetesFleetManagerDataSourceModel struct {
	Location          string                   `tfschema:"location"`
	Name              string                   `tfschema:"name"`
	ResourceGroupName string                   `tfschema:"resource_group_name"`
	HubProfile        []FleetManagerHubProfile `tfschema:"hub_profile"`
	Tags              map[string]interface{}   `tfschema:"tags"`
}

func (KubernetesFleetManagerDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (KubernetesFleetManagerDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"hub_profile": {
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"agent_profile": {
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"subnet_id": {
									Computed: true,
									Type:     pluginsdk.TypeString,
								},
								"virtual_machine_size": {
									Computed: true,
									Type:     pluginsdk.TypeString,
								},
							},
						},
						Type: pluginsdk.TypeList,
					},
					"api_server_access_profile": {
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enable_private_cluster": {
									Computed: true,
									Type:     pluginsdk.TypeBool,
								},
							},
						},
						Type: pluginsdk.TypeList,
					},
					"dns_prefix": {
						Computed: true,
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
					"portal_fqdn": {
						Computed: true,
						Type:     pluginsdk.TypeString,
					},
				},
			},
			Type: pluginsdk.TypeList,
		},
		"location": commonschema.LocationComputed(),
		"tags":     commonschema.TagsDataSource(),
	}
}

func (KubernetesFleetManagerDataSource) ModelObject() interface{} {
	return &KubernetesFleetManagerDataSourceModel{}
}

func (KubernetesFleetManagerDataSource) ResourceType() string {
	return "azurerm_kubernetes_fleet_manager"
}

func (KubernetesFleetManagerDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerService.V20240401.Fleets
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state KubernetesFleetManagerDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := fleets.NewFleetID(subscriptionId, state.ResourceGroupName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.SetID(id)

			if model := resp.Model; model != nil {
				if err := mapFleetToKubernetesFleetManagerDataSourceModel(*model, &state); err != nil {
					return err
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func mapFleetToKubernetesFleetManagerDataSourceModel(input fleets.Fleet, output *KubernetesFleetManagerDataSourceModel) error {
	output.Location = location.Normalize(input.Location)
	output.Tags = tags.Flatten(input.Tags)

	if input.Properties == nil {
		input.Properties = &fleets.FleetProperties{}
	}

	hubProfile, err := flattenFleetManagerHubProfile(input.Properties.HubProfile)
	if err != nil {
		return err
	}
	output.HubProfile = hubProfile
	return nil
}
