package netapp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2021-10-01/volumegroups"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	netAppValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetAppVolumeGroupResource struct{}

type NetAppVolumeGroupModel struct {
	Name                  string                    `tfschema:"name"`
	ResourceGroupName     string                    `tfschema:"resource_group_name"`
	Location              string                    `tfschema:"location"`
	AccountName           string                    `tfschema:"account_name"`
	GroupDescription      string                    `tfschema:"group_description"`
	ApplicationType       string                    `tfschema:"application_type"`
	ApplicationIdentifier string                    `tfschema:"application_identifier"`
	DeploymentSpecId      string                    `tfschema:"deployment_spec_id"`
	Tags                  map[string]interface{}    `tfschema:"tags"`
	Volumes               []NetAppVolumeGroupVolume `tfschema:"volume"`
}

var _ sdk.Resource = NetAppVolumeGroupResource{}

func (r NetAppVolumeGroupResource) ModelObject() interface{} {
	return &NetAppVolumeGroupModel{}
}

func (r NetAppVolumeGroupResource) ResourceType() string {
	return "azurerm_netapp_volume_group"
}

func (r NetAppVolumeGroupResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return volumegroups.ValidateVolumeGroupID
}

func (r NetAppVolumeGroupResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": azure.SchemaResourceGroupName(),

		"location": azure.SchemaLocation(),

		"account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: netAppValidate.AccountName,
		},

		"group_description": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"application_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(volumegroups.ApplicationTypeSAPNegativeHANA),
			}, false),
		},

		"application_identifier": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 3),
		},

		"deployment_spec_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"volume": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MinItems: 5,
			MaxItems: 5,
			Elem: &pluginsdk.Resource{
				Schema: netAppVolumeGroupVolumeSchema(),
			},
		},

		// Can't use tags.Schema since there is no patch available
		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		//"tags": commonschema.Tags(),
	}
}

func (r NetAppVolumeGroupResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		/*
			TODO - This section is for `Computed: true` only items, i.e. useful values that are returned by the
			datasource that can be used as outputs or passed programmatically to other resources or data sources.

			TODO (pmarques) - use this for first level attributes when Volume resource gets migrated to tfschema
		*/
	}
}

func (r NetAppVolumeGroupResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.VolumeGroupClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model NetAppVolumeGroupModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := volumegroups.NewVolumeGroupID(subscriptionId, model.ResourceGroupName, model.AccountName, model.Name)

			metadata.Logger.Infof("Import check for %s", id)
			existing, err := client.VolumeGroupsGet(ctx, id)
			if err != nil && existing.HttpResponse.StatusCode != http.StatusNotFound {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if existing.Model != nil && existing.Model.Id != nil && *existing.Model.Id != "" {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			applicationType := volumegroups.ApplicationType(model.ApplicationType)

			volumeList, err := expandNetAppVolumeGroupVolumes(model.Volumes, id)
			if err != nil {
				return err
			}

			parameters := volumegroups.VolumeGroupDetails{
				Location: utils.String(location.Normalize(model.Location)),
				Properties: &volumegroups.VolumeGroupProperties{
					GroupMetaData: &volumegroups.VolumeGroupMetaData{
						GroupDescription:      utils.String(model.GroupDescription),
						ApplicationType:       &applicationType,
						ApplicationIdentifier: utils.String(model.ApplicationIdentifier),
						DeploymentSpecId:      utils.String(model.DeploymentSpecId),
					},
					Volumes: volumeList,
				},
				Tags: tags.Expand(model.Tags),
			}

			err = client.VolumeGroupsCreateThenPoll(ctx, id, parameters)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			// TODO: Check if this is necessary for volume groups
			// // Waiting for volume be completely provisioned
			// if err := waitForVolumeCreateOrUpdate(ctx, client, id); err != nil {
			// 	return err
			// }

			metadata.SetID(id)

			return nil
		},
	}
}

func (r NetAppVolumeGroupResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {

			client := metadata.Client.NetApp.VolumeGroupClient

			id, err := volumegroups.ParseVolumeGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Decoding state for %s", id)
			var state NetAppVolumeGroupModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			existing, err := client.VolumeGroupsGet(ctx, *id)
			if err != nil {
				if existing.HttpResponse.StatusCode == http.StatusNotFound {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			model := NetAppVolumeGroupModel{
				Name:              id.VolumeGroupName,
				AccountName:       id.AccountName,
				Location:          location.NormalizeNilable(existing.Model.Location),
				ResourceGroupName: id.ResourceGroupName,
				Tags:              tags.Flatten(existing.Model.Tags),
			}

			if props := existing.Model.Properties; props != nil {
				model.GroupDescription = utils.NormalizeNilableString(props.GroupMetaData.GroupDescription)
				model.ApplicationIdentifier = utils.NormalizeNilableString(props.GroupMetaData.ApplicationIdentifier)
				model.DeploymentSpecId = utils.NormalizeNilableString(props.GroupMetaData.DeploymentSpecId)

				volumes, err := flattenNetAppVolumeGroupVolumes(props.Volumes)
				if err != nil {
					return fmt.Errorf("setting `volume`: %+v", err)
				}

				model.Volumes = volumes
			}

			return metadata.Encode(&model)
		},
	}
}

func (r NetAppVolumeGroupResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// TODO - Delete Func
			return nil
		},
	}
}
