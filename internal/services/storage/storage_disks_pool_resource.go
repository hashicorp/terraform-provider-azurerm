package storage

import (
	"context"
	"fmt"
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
	"regexp"
	"time"
)

type DisksPoolResource struct{}

var _ sdk.ResourceWithUpdate = DisksPoolResource{}

type DiskPoolJobModel struct {
	AdditionalCapabilities []string               `tfschema:"additional_capabilities"` // List of additional capabilities for a Disk Pool.
	AvailabilityZones      []string               `tfschema:"availability_zones"`      // Logical zone for Disk Pool resource; example: [\"1\"].
	Location               string                 `tfschema:"location"`                // The geo-location where the resource lives.
	ManagedBy              string                 `tfschema:"managed_by"`              // Azure resource id. Indicates if this resource is managed by another Azure resource.
	ManagedByExtended      []string               `tfschema:"managed_by_extended"`     // List of Azure resource ids that manage this resource.
	Name                   string                 `tfschema:"name"`                    // The name of the Disk Pool.
	ResourceGroupName      string                 `tfschema:"resource_group_name"`     // The name of the resource group. The name is case insensitive.
	Sku                    []DiskPoolSku          `tfschema:"sku"`                     // Determines the SKU of the Disk Pool
	SubnetId               string                 `tfschema:"subnet_id"`               // Azure Resource ID of a Subnet for the Disk Pool.
	Tags                   map[string]interface{} `tfschema:"tags"`                    // Resource tags.
}

type DiskPoolSku struct {
	Name string `tfschema:"name"` // Determines the SKU of the Disk Pool
	// Tier - Sku tier
	Tier string `tfschema:"tier"` // Sku tier
}

func (d DisksPoolResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"additional_capabilities": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
		"availability_zones": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
		"location": location.Schema(),
		"managed_by_extended": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: azure.ValidateResourceID,
			},
		},
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.All(
				validation.StringIsNotEmpty,
				validation.StringLenBetween(7, 30),
				validation.StringMatch(
					regexp.MustCompile("^[A-Za-z\\d][A-Za-z\\d.\\-_]*[A-Za-z\\d_]$"),
					"The name must begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods, or hyphens.",
				),
			),
		},
		"resource_group_name": azure.SchemaResourceGroupName(),
		"sku": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"Basic_B1",
							"Standard_S1",
							"Premium_P1",
						}, false),
					},
					"tier": {
						Type:         pluginsdk.TypeString,
						Computed:     true,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(possibleDiskPoolTierValues(), false),
					},
				},
			},
		},
		"subnet_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: networkValidate.SubnetID,
		},
		"tags": {
			Type:         pluginsdk.TypeMap,
			Optional:     true,
			ValidateFunc: tags.Validate,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (d DisksPoolResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"managed_by": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (d DisksPoolResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			diskPool := DiskPoolJobModel{}
			err := metadata.Decode(&diskPool)
			if err != nil {
				return err
			}
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := parse.NewStorageDisksPoolID(subscriptionId, diskPool.ResourceGroupName, diskPool.Name)

			client := metadata.Client.Storage.DisksPoolsClient

			if metadata.ResourceData.IsNewResource() {
				existing, err := client.Get(ctx, diskPool.ResourceGroupName, diskPool.Name)
				notExistingResp := utils.ResponseWasNotFound(existing.Response)
				if err != nil && !notExistingResp {
					return fmt.Errorf("checking for presence of existing %q: %+v", id, err)
				}
				if !notExistingResp {
					return metadata.ResourceRequiresImport(d.ResourceType(), id)
				}
			}

			createParameter := storagepool.DiskPoolCreate{
				DiskPoolCreateProperties: &storagepool.DiskPoolCreateProperties{
					AvailabilityZones:      &diskPool.AvailabilityZones,
					SubnetID:               utils.String(diskPool.SubnetId),
					AdditionalCapabilities: &diskPool.AdditionalCapabilities,
				},
				Location:          utils.String(diskPool.Location),
				ManagedBy:         utils.String(diskPool.ManagedBy),
				ManagedByExtended: &diskPool.ManagedByExtended,
				Name:              utils.String(diskPool.Name),
				Sku:               expandDisksPoolSku(diskPool),
				Tags:              tags.Expand(diskPool.Tags),
			}
			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, diskPool.Name, createParameter)
			if err != nil {
				return err
			}
			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of DisksPool %q (Resource Group %q): %+v", diskPool.Name, diskPool.ResourceGroupName, err)
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
			diskPoolId := parse.NewStorageDisksPoolID(id.SubscriptionId, id.ResourceGroup, id.DiskPoolName)
			client := metadata.Client.Storage.DisksPoolsClient
			resp, err := client.Get(ctx, diskPoolId.ResourceGroup, diskPoolId.DiskPoolName)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			model := DiskPoolJobModel{
				AdditionalCapabilities: nil,
				AvailabilityZones:      *resp.AvailabilityZones,
				Location:               *resp.Location,
				ManagedBy:              "",
				ManagedByExtended:      nil,
				Name:                   *resp.Name,
				ResourceGroupName:      id.ResourceGroup,
				Sku:                    flattenDiskPoolSku(*resp.Sku),
				SubnetId:               *resp.SubnetID,
				Tags:                   tags.Flatten(resp.Tags),
			}
			if resp.AdditionalCapabilities != nil {
				model.AdditionalCapabilities = *resp.AdditionalCapabilities
			}
			if resp.ManagedBy != nil {
				model.ManagedBy = *resp.ManagedBy
			}
			if resp.ManagedByExtended != nil {
				model.ManagedByExtended = *resp.ManagedByExtended
			}
			return metadata.Encode(&model)
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
			locks.ByID(metadata.ResourceData.Id())
			defer locks.UnlockByID(metadata.ResourceData.Id())
			poolId := parse.NewStorageDisksPoolID(id.SubscriptionId, id.ResourceGroup, id.DiskPoolName)
			client := metadata.Client.Storage.DisksPoolsClient
			future, err := client.Delete(ctx, poolId.ResourceGroup, poolId.DiskPoolName)
			if err != nil {
				return err
			}
			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for deletion of DiskPool %q (Resource Group %q): %+v", poolId.DiskPoolName, poolId.ResourceGroup, err)
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
			poolId := parse.NewStorageDisksPoolID(id.SubscriptionId, id.ResourceGroup, id.DiskPoolName)
			client := metadata.Client.Storage.DisksPoolsClient
			patch := storagepool.DiskPoolUpdate{
				ManagedBy: nil,
			}

			m := DiskPoolJobModel{}
			err = metadata.Decode(&m)
			if err != nil {
				return err
			}
			if r.HasChange("managed_by_extended") {
				patch.ManagedByExtended = &m.ManagedByExtended
			}
			if r.HasChange("sku") {
				patch.Sku = &storagepool.Sku{
					Name: utils.String(m.Sku[0].Name),
					Tier: utils.String(m.Sku[0].Tier),
				}
			}
			if r.HasChange("tags") {
				patch.Tags = tags.Expand(m.Tags)
			}
			future, err := client.Update(ctx, poolId.ResourceGroup, poolId.DiskPoolName, patch)
			if err != nil {
				return err
			}
			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for update of DiskPool %q (Resource Group %q): %+v", poolId.DiskPoolName, poolId.ResourceGroup, err)
			}
			return nil
		},
	}
}

func (d DisksPoolResource) ModelObject() interface{} {
	return &DiskPoolJobModel{}
}

func (d DisksPoolResource) ResourceType() string {
	return "azurerm_storage_disks_pool"
}

func expandDisksPoolSku(diskPool DiskPoolJobModel) *storagepool.Sku {
	return &storagepool.Sku{
		Name: utils.String((diskPool.Sku)[0].Name),
		Tier: utils.String((diskPool.Sku)[0].Tier),
	}
}

func flattenDiskPoolSku(sku storagepool.Sku) []DiskPoolSku {
	r := DiskPoolSku{
		Name: *sku.Name,
	}
	if sku.Tier != nil {
		r.Tier = *sku.Tier
	}
	return []DiskPoolSku{r}
}

func possibleDiskPoolTierValues() []string {
	var tiers []string
	for _, tier := range storagepool.PossibleDiskPoolTierValues() {
		tiers = append(tiers, string(tier))
	}
	return tiers
}
