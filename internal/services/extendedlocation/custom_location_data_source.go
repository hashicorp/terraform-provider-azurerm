package extendedlocation

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
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type CustomLocationDataSource struct{}

var _ sdk.DataSource = CustomLocationDataSource{}

type CustomLocationModel struct {
	ClusterExtensionIds []string                       `tfschema:"cluster_extension_ids"`
	DisplayName         string                         `tfschema:"display_name"`
	HostResourceId      string                         `tfschema:"host_resource_id"`
	Identities          []identity.ModelSystemAssigned `tfschema:"identities"`
	Location            string                         `tfschema:"location"`
	Name                string                         `tfschema:"name"`
	Namespace           string                         `tfschema:"namespace"`
	ResourceGroupName   string                         `tfschema:"resource_group_name"`
	Tags                map[string]interface{}         `tfschema:"tags"`
}

func (r CustomLocationDataSource) ResourceType() string {
	return "azurerm_custom_location"
}

func (r CustomLocationDataSource) ModelObject() interface{} {
	return &CustomLocationModel{}
}

func (r CustomLocationDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (r CustomLocationDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"cluster_extension_ids": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
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

		"identities": commonschema.SystemAssignedIdentityComputed(),

		"location": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"namespace": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (r CustomLocationDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ExtendedLocation.CustomLocations
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model CustomLocationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := customlocations.NewCustomLocationID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s does not exist", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := CustomLocationModel{
				Name:              model.Name,
				ResourceGroupName: model.ResourceGroupName,
			}

			if model := existing.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = tags.Flatten(model.Tags)
				state.Identities = identity.FlattenSystemAssignedToModel(model.Identity)

				if props := model.Properties; props != nil {
					state.ClusterExtensionIds = pointer.From(props.ClusterExtensionIds)
					state.DisplayName = pointer.From(props.DisplayName)
					state.HostResourceId = pointer.From(props.HostResourceId)
					state.Namespace = pointer.From(props.Namespace)
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
