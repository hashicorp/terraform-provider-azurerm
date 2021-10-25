package storage

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/storagepool/mgmt/2021-08-01/storagepool"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"regexp"
)

type DiskPoolResource struct{}

var _ sdk.ResourceWithUpdate = DiskPoolResource{}

type DiskPoolJobModel struct {
	AdditionalCapabilities []string `tfschema:"additional_capabilities"` //List of additional capabilities for a Disk Pool.
	AvailabilityZones      []string `tfschema:"availability_zones"`      //Logical zone for Disk Pool resource; example: [\"1\"].
	Location               string   `tfschema:"location"`                //The geo-location where the resource lives.
	//TODO: Figure out how MangedBy and MangedByExtended work
	ManagedBy         string                 `tfschema:"managed_by"`          //Azure resource id. Indicates if this resource is managed by another Azure resource.
	ManagedByExtended []string               `tfschema:"managed_by_extended"` //List of Azure resource ids that manage this resource.
	Name              string                 `tfschema:"name"`                //The name of the Disk Pool.
	ResourceGroupName string                 `tfschema:"resource_group_name"` //The name of the resource group. The name is case insensitive.
	Sku               *DiskPoolSku           `tfschema:"sku"`                 //Determines the SKU of the Disk Pool
	SubnetId          string                 `tfschema:"subnetId"`            //Azure Resource ID of a Subnet for the Disk Pool.
	Tags              map[string]interface{} `tfschema:"tags"`                //Resource tags.
	Type              string                 `tfschema:"type"`                // Ex- Microsoft.Compute/virtualMachines or Microsoft.Storage/storageAccounts
}

type DiskPoolSku struct {
	Name string `tfschema:"name"` //Determines the SKU of the Disk Pool
	// Tier - Sku tier
	Tier string `tfschema:"tier"` //Sku tier
}

func (d DiskPoolResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"additional_capabilities": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
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
		"disk_ids": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: computeValidate.ManagedDiskID,
			},
		},
		"location": location.Schema(),
		"managed_by": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: azure.ValidateResourceID,
		},
		"managed_by_extended": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Elem: pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: azure.ValidateResourceID,
			},
		},
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringIsNotEmpty,
				validation.StringLenBetween(7, 30),
				validation.StringMatch(
					regexp.MustCompile("^[A-Za-z\\d][A-Za-z\\d\\.\\-_]*[A-Za-z\\d_]$"),
					"The name must begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods, or hyphens.",
				),
			),
		},
		"resource_group_name": azure.SchemaResourceGroupName(),
		"sku": {
			Type:     pluginsdk.TypeList,
			MinItems: 1,
			MaxItems: 1,
			Required: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice([]string{"B1", "S1", "P1"}, false),
					},
					"tier": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"Basic", "Standard", "Premium"}, false),
					},
				},
			},
		},
		"subnetId": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: networkValidate.SubnetID,
		},
		"type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (d DiskPoolResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			diskPool := DiskPoolJobModel{}
			err := metadata.Decode(&diskPool)
			if err != nil {
				return err
			}
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := parse.NewDiskPoolID(subscriptionId, diskPool.ResourceGroupName, diskPool.Name)
			client := metadata.Client.Storage.DiskPoolsClient

			//TODO: Finish after read was complete
			if metadata.ResourceData.IsNewResource() {
				//existing, err := r.getJob(ctx, client, id)
				//if err != nil {
				//	if !utils.ResponseWasNotFound(existing.Response) {
				//		return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				//	}
				//}
				//if !utils.ResponseWasNotFound(existing.Response) {
				//	return metadata.ResourceRequiresImport(r.ResourceType(), id)
				//}
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
				Sku: &storagepool.Sku{
					Name: utils.String(diskPool.Sku.Name),
					Tier: utils.String(diskPool.Sku.Tier),
				},
				Tags: tags.Expand(diskPool.Tags),
				Type: utils.String(diskPool.Type),
			}
			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, diskPool.Name, createParameter)
			if err != nil {
				return err
			}
			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of DiskPool %q (Resource Group %q): %+v", diskPool.Name, diskPool.ResourceGroupName, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (d DiskPoolResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := parse.DiskPoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			diskPoolId := parse.NewDiskPoolID(id.SubscriptionId, id.ResourceGroup, id.Name)
			client := metadata.Client.Storage.DiskPoolsClient
			resp, err := client.Get(ctx, diskPoolId.ResourceGroup, diskPoolId.Name)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			model := DiskPoolJobModel{
				AvailabilityZones: *resp.AvailabilityZones,
				Name:              *resp.Name,
				ResourceGroupName: id.ResourceGroup,
				Location:          *resp.Location,
				SubnetId:          *resp.SubnetID,
				Sku:               flattenDiskPoolSku(*resp.Sku),
				Tags:              tags.Flatten(resp.Tags),
			}
			if resp.Type != nil {
				model.Type = *resp.Type
			}
			if resp.ManagedBy != nil {
				model.ManagedBy = *resp.ManagedBy
			}
			if resp.ManagedByExtended != nil {
				model.ManagedByExtended = *resp.ManagedByExtended
			}

			return nil
		},
	}
}

func (d DiskPoolResource) Delete() sdk.ResourceFunc {
	panic("implement me")
}

func (d DiskPoolResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	panic("implement me")
}

func (d DiskPoolResource) Update() sdk.ResourceFunc {
	panic("implement me")
}

func (d DiskPoolResource) Attributes() map[string]*schema.Schema {
	panic("implement me")
}

func (d DiskPoolResource) ModelObject() interface{} {
	panic("implement me")
}

func (d DiskPoolResource) ResourceType() string {
	panic("implement me")
}
func flattenDiskPoolSku(sku storagepool.Sku) *DiskPoolSku {
	r := &DiskPoolSku{
		Name: *sku.Name,
	}
	if sku.Tier != nil {
		r.Tier = *sku.Tier
	}
	return r
}
