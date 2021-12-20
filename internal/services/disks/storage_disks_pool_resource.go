package disks

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/disks/sdk/2021-08-01/diskpools"
	disksValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/disks/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StorageDisksPoolResource struct{}

var _ sdk.ResourceWithUpdate = StorageDisksPoolResource{}

type StorageDisksPoolJobModel struct {
	Name              string                 `tfschema:"name"`
	ResourceGroupName string                 `tfschema:"resource_group_name"`
	Location          string                 `tfschema:"location"`
	AvailabilityZones []string               `tfschema:"availability_zones"`
	Sku               string                 `tfschema:"sku_name"`
	SubnetId          string                 `tfschema:"subnet_id"`
	Tags              map[string]interface{} `tfschema:"tags"`
}

func (d StorageDisksPoolResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: disksValidate.DiskPoolName(),
		},
		"resource_group_name": commonschema.ResourceGroupName(),
		"location":            commonschema.Location(),
		"availability_zones": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
		"sku_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: disksValidate.DiskPoolSku(),
		},
		"subnet_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: networkValidate.SubnetID,
		},
		"tags": commonschema.Tags(),
	}
}

func (d StorageDisksPoolResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (d StorageDisksPoolResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.Disks.DiskPoolsClient

			m := StorageDisksPoolJobModel{}
			err := metadata.Decode(&m)
			if err != nil {
				return err
			}

			id := diskpools.NewDiskPoolID(subscriptionId, m.ResourceGroupName, m.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %q: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(d.ResourceType(), id)
			}

			createParameter := diskpools.DiskPoolCreate{
				Name:     utils.String(m.Name),
				Location: location.Normalize(m.Location),
				Properties: diskpools.DiskPoolCreateProperties{
					AvailabilityZones: &m.AvailabilityZones,
					SubnetId:          m.SubnetId,
				},
				Sku:  expandDisksPoolSku(m.Sku),
				Tags: tags.Expand(m.Tags),
			}
			if err := client.CreateOrUpdateThenPoll(ctx, id, createParameter); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)
			return nil
		},
	}
}

func (d StorageDisksPoolResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Disks.DiskPoolsClient
			id, err := diskpools.ParseDiskPoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			m := StorageDisksPoolJobModel{
				Name:              id.DiskPoolName,
				ResourceGroupName: id.ResourceGroupName,
			}
			if model := resp.Model; model != nil {
				if model.Sku != nil {
					m.Sku = model.Sku.Name
				}
				m.Tags = flattenTags(model.Tags)

				m.AvailabilityZones = model.Properties.AvailabilityZones
				m.Location = location.Normalize(model.Location)
				m.SubnetId = model.Properties.SubnetId
			}

			return metadata.Encode(&m)
		},
	}
}

func (d StorageDisksPoolResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Disks.DiskPoolsClient
			id, err := diskpools.ParseDiskPoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (d StorageDisksPoolResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.StorageDisksPoolID
}

func (d StorageDisksPoolResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Disks.DiskPoolsClient
			id, err := diskpools.ParseDiskPoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(metadata.ResourceData.Id())
			defer locks.UnlockByID(metadata.ResourceData.Id())

			patch := diskpools.DiskPoolUpdate{}
			var m StorageDisksPoolJobModel
			if err = metadata.Decode(&m); err != nil {
				return fmt.Errorf("decoding model: %+v", err)
			}

			if metadata.ResourceData.HasChange("sku") {
				sku := expandDisksPoolSku(m.Sku)
				patch.Sku = &sku
			}
			if metadata.ResourceData.HasChange("tags") {
				patch.Tags = tags.Expand(m.Tags)
			}

			if err := client.UpdateThenPoll(ctx, *id, patch); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (d StorageDisksPoolResource) ModelObject() interface{} {
	return &StorageDisksPoolJobModel{}
}

func (d StorageDisksPoolResource) ResourceType() string {
	return "azurerm_storage_disks_pool"
}

func expandDisksPoolSku(sku string) diskpools.Sku {
	parts := strings.Split(sku, "_")
	return diskpools.Sku{
		Name: sku,
		Tier: &parts[0],
	}
}

func flattenTags(input *map[string]string) map[string]interface{} {
	output := make(map[string]interface{})

	if input != nil {
		for k, v := range *input {
			val := v
			output[k] = val
		}
	}

	return output
}
