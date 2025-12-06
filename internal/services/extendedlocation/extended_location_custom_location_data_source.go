// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package extendedlocation

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.DataSource = ExtendedLocationCustomLocationDataSource{}

type ExtendedLocationCustomLocationDataSource struct{}

type ExtendedLocationCustomLocationDataSourceModel struct {
	Authentication      []AuthModel `tfschema:"authentication"`
	ClusterExtensionIds []string    `tfschema:"cluster_extension_ids"`
	DisplayName         string      `tfschema:"display_name"`
	HostResourceId      string      `tfschema:"host_resource_id"`
	HostType            string      `tfschema:"host_type"`
	Location            string      `tfschema:"location"`
	Name                string      `tfschema:"name"`
	Namespace           string      `tfschema:"namespace"`
	ResourceGroupName   string      `tfschema:"resource_group_name"`
}

func (r ExtendedLocationCustomLocationDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[A-Za-z\d.\-_]*[A-Za-z\d]$`),
				"supported alphanumeric characters and periods, underscores, hyphens. Name should end with an alphanumeric character.",
			),
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (r ExtendedLocationCustomLocationDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*schema.Schema{
		"authentication": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"value": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"cluster_extension_ids": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"host_resource_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"host_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"location": commonschema.LocationComputed(),

		"namespace": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r ExtendedLocationCustomLocationDataSource) ModelObject() interface{} {
	return &ExtendedLocationCustomLocationDataSourceModel{}
}

func (r ExtendedLocationCustomLocationDataSource) ResourceType() string {
	return "azurerm_extended_location_custom_location"
}

func (r ExtendedLocationCustomLocationDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.ExtendedLocation.CustomLocationsClient

			var model ExtendedLocationCustomLocationResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := customlocations.NewCustomLocationID(subscriptionId, model.ResourceGroupName, model.Name)
			resp, err := client.Get(ctx, id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := ExtendedLocationCustomLocationDataSourceModel{
				Name:              id.CustomLocationName,
				ResourceGroupName: id.ResourceGroupName,
			}
			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)

				if props := model.Properties; props != nil {
					state.ClusterExtensionIds = pointer.From(props.ClusterExtensionIds)
					state.DisplayName = pointer.From(props.DisplayName)
					state.HostResourceId = pointer.From(props.HostResourceId)
					state.HostType = string(pointer.From(props.HostType))
					state.Namespace = pointer.From(props.Namespace)

					if props.Authentication != nil && props.Authentication.Type != nil && props.Authentication.Value != nil {
						state.Authentication = []AuthModel{
							{
								Type:  pointer.From(props.Authentication.Type),
								Value: pointer.From(props.Authentication.Value),
							},
						}
					}
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
