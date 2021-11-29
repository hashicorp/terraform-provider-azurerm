package storage

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storagepool/mgmt/2021-08-01/storagepool"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DisksPoolResource struct{}

var _ sdk.ResourceWithUpdate = DisksPoolResource{}

type DisksPoolJobModel struct {
	Name              string                 `tfschema:"name"`
	ResourceGroupName string                 `tfschema:"resource_group_name"`
	Location          string                 `tfschema:"location"`
	AvailabilityZones []string               `tfschema:"availability_zones"`
	Sku               string                 `tfschema:"sku_name"`
	SubnetId          string                 `tfschema:"subnet_id"`
	Tags              map[string]interface{} `tfschema:"tags"`
}

func (d DisksPoolResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringIsNotEmpty,
				validation.StringLenBetween(7, 30),
				validation.StringMatch(
					regexp.MustCompile(`^[A-Za-z\d][A-Za-z\d.\-_]*[A-Za-z\d_]$`),
					"The name must begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods, or hyphens.",
				),
			),
		},
		"resource_group_name": azure.SchemaResourceGroupName(),
		"location":            location.Schema(),
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
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice(
				[]string{
					"Basic_B1",
					"Standard_S1",
					"Premium_P1",
				}, false,
			),
		},
		"subnet_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: networkValidate.SubnetID,
		},
		"tags": tags.Schema(),
	}
}

func (d DisksPoolResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (d DisksPoolResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			m := DisksPoolJobModel{}
			err := metadata.Decode(&m)
			if err != nil {
				return err
			}
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := parse.NewStorageDisksPoolID(subscriptionId, m.ResourceGroupName, m.Name)

			client := metadata.Client.Storage.DisksPoolsClient

			existing, err := client.Get(ctx, m.ResourceGroupName, m.Name)
			notExistingResp := utils.ResponseWasNotFound(existing.Response)
			if err != nil && !notExistingResp {
				return fmt.Errorf("checking for presence of existing %q: %+v", id, err)
			}
			if !notExistingResp {
				return metadata.ResourceRequiresImport(d.ResourceType(), id)
			}

			createParameter := storagepool.DiskPoolCreate{
				DiskPoolCreateProperties: &storagepool.DiskPoolCreateProperties{
					AvailabilityZones: &m.AvailabilityZones,
					SubnetID:          &m.SubnetId,
				},
				Location: utils.String(m.Location),
				Name:     utils.String(m.Name),
				Sku:      expandDisksPoolSku(m.Sku),
				Tags:     tags.Expand(m.Tags),
			}
			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, m.Name, createParameter)
			if err != nil {
				return fmt.Errorf("creation of %q: %+v", id, err)
			}
			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of %q: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (d DisksPoolResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := parse.StorageDisksPoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			client := metadata.Client.Storage.DisksPoolsClient
			resp, err := client.Get(ctx, id.ResourceGroup, id.DiskPoolName)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %q: %+v", id, err)
			}
			m := DisksPoolJobModel{
				Name:              id.DiskPoolName,
				ResourceGroupName: id.ResourceGroup,
				Tags:              tags.Flatten(resp.Tags),
			}
			if resp.AvailabilityZones != nil {
				m.AvailabilityZones = *resp.AvailabilityZones
			}
			if resp.Location != nil {
				m.Location = location.Normalize(*resp.Location)
			}
			if resp.Sku != nil && resp.Sku.Name != nil {
				m.Sku = *resp.Sku.Name
			}
			if resp.SubnetID != nil {
				m.SubnetId = *resp.SubnetID
			}
			return metadata.Encode(&m)
		},
	}
}

func (d DisksPoolResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := parse.StorageDisksPoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			client := metadata.Client.Storage.DisksPoolsClient
			future, err := client.Delete(ctx, id.ResourceGroup, id.DiskPoolName)
			if err != nil {
				return fmt.Errorf("deletion of %q: %+v", id, err)
			}
			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for deletion of %q : %+v", id, err)
			}
			return nil
		},
	}
}

func (d DisksPoolResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.StorageDisksPoolID
}

func (d DisksPoolResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			r := metadata.ResourceData
			id, err := parse.StorageDisksPoolID(r.Id())
			if err != nil {
				return err
			}

			locks.ByID(r.Id())
			defer locks.UnlockByID(r.Id())

			client := metadata.Client.Storage.DisksPoolsClient
			patch := storagepool.DiskPoolUpdate{}
			m := DisksPoolJobModel{}
			err = metadata.Decode(&m)
			if err != nil {
				return err
			}

			if r.HasChange("sku") {
				patch.Sku = expandDisksPoolSku(m.Sku)
			}
			if r.HasChange("tags") {
				patch.Tags = tags.Expand(m.Tags)
			}

			future, err := client.Update(ctx, id.ResourceGroup, id.DiskPoolName, patch)
			if err != nil {
				return fmt.Errorf("update of %q: %+v", id, err)
			}
			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for update of %q : %+v", id, err)
			}
			return nil
		},
	}
}

func (d DisksPoolResource) ModelObject() interface{} {
	return &DisksPoolJobModel{}
}

func (d DisksPoolResource) ResourceType() string {
	return "azurerm_storage_disks_pool"
}

func expandDisksPoolSku(sku string) *storagepool.Sku {
	parts := strings.Split(sku, "_")
	return &storagepool.Sku{
		Name: &sku,
		Tier: &parts[0],
	}
}
