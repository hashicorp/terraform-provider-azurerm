package extendedlocation

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.DataSource = CustomLocationDataSource{}

type CustomLocationDataSource struct{}

type CustomLocationDataSourceModel struct {
	Name                string      `tfschema:"name"`
	ResourceGroupName   string      `tfschema:"resource_group_name"`
	Location            string      `tfschema:"location"`
	Authentication      []AuthModel `tfschema:"authentication"`
	ClusterExtensionIds []string    `tfschema:"cluster_extension_ids"`
	DisplayName         string      `tfschema:"display_name"`
	HostResourceId      string      `tfschema:"host_resource_id"`
	HostType            string      `tfschema:"host_type"`
	Namespace           string      `tfschema:"namespace"`
}

func (r CustomLocationDataSource) Arguments() map[string]*pluginsdk.Schema {
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

func (r CustomLocationDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*schema.Schema{
		"location": commonschema.LocationComputed(),

		"namespace": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"host_resource_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"cluster_extension_ids": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"host_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

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
	}
}

func (r CustomLocationDataSource) ModelObject() interface{} {
	return &CustomLocationDataSourceModel{}
}

func (r CustomLocationDataSource) ResourceType() string {
	return "azurerm_extended_location_custom_location"
}

func (r CustomLocationDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model CustomLocationResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.ExtendedLocation.CustomLocationsClient

			id := customlocations.NewCustomLocationID(subscriptionId, model.ResourceGroupName, model.Name)
			resp, err := client.Get(ctx, id)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", id, err)
			}

			state := CustomLocationDataSourceModel{
				Name:              id.CustomLocationName,
				ResourceGroupName: id.ResourceGroupName,
			}
			if model := resp.Model; model != nil {
				state.Location = model.Location

				if props := model.Properties; props != nil {
					state.ClusterExtensionIds = pointer.From(props.ClusterExtensionIds)
					state.DisplayName = pointer.From(props.DisplayName)
					state.HostResourceId = pointer.From(props.HostResourceId)
					state.HostType = string(pointer.From(props.HostType))
					state.Namespace = pointer.From(props.Namespace)

					// API always returns an empty `authentication` block even it's not specified. Tracing the bug: https://github.com/Azure/azure-rest-api-specs/issues/30101
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
