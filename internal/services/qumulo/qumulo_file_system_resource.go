package qumulo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/subnets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/qumulostorage/2022-10-12/filesystems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/qumulo/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.Resource           = FileSystemResource{}
	_ sdk.ResourceWithUpdate = FileSystemResource{}
)

const (
	offerId     = "qumulo-saas-mpp"
	publisherId = "qumulo1584033880660"
)

type FileSystemResource struct{}

func (r FileSystemResource) ModelObject() interface{} {
	return &FileSystemResourceSchema{}
}

type FileSystemResourceSchema struct {
	AdminPassword     string                 `tfschema:"admin_password"`
	InitialCapacity   int64                  `tfschema:"initial_capacity"`
	Location          string                 `tfschema:"location"`
	MarketplacePlanId string                 `tfschema:"marketplace_plan_id"`
	Name              string                 `tfschema:"name"`
	ResourceGroupName string                 `tfschema:"resource_group_name"`
	StorageSku        string                 `tfschema:"storage_sku"`
	SubnetId          string                 `tfschema:"subnet_id"`
	Tags              map[string]interface{} `tfschema:"tags"`
	UserEmailAddress  string                 `tfschema:"user_email_address"`
	Zone              string                 `tfschema:"zone"`
}

func (r FileSystemResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return filesystems.ValidateFileSystemID
}

func (r FileSystemResource) ResourceType() string {
	return "azurerm_qumulo_file_system"
}

func (r FileSystemResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			ForceNew:     true,
			Required:     true,
			Type:         pluginsdk.TypeString,
			ValidateFunc: validate.FileSystemName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"admin_password": {
			Type:      pluginsdk.TypeString,
			Required:  true,
			ForceNew:  true,
			Sensitive: true,
		},

		"initial_capacity": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(18, 1000),
		},

		"marketplace_plan_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"storage_sku": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(filesystems.PossibleValuesForStorageSku(), false),
		},

		"subnet_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"user_email_address": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"zone": commonschema.ZoneSingleOptionalForceNew(),

		"location": commonschema.Location(),

		"tags": commonschema.Tags(),
	}
}

func (r FileSystemResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r FileSystemResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Qumulo.FileSystemsClient

			var config FileSystemResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId

			id := filesystems.NewFileSystemID(subscriptionId, config.ResourceGroupName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			// check the subnet is valid: https://learn.microsoft.com/en-us/azure/partner-solutions/qumulo/qumulo-troubleshoot#you-cant-create-a-resource
			if err := checkSubnet(ctx, config.SubnetId, metadata); err != nil {
				return err
			}

			payload := filesystems.FileSystemResource{
				Location: location.Normalize(config.Location),
				Tags:     tags.Expand(config.Tags),
				Properties: filesystems.FileSystemResourceProperties{
					AdminPassword:     config.AdminPassword,
					AvailabilityZone:  pointer.To(config.Zone),
					DelegatedSubnetId: config.SubnetId,
					InitialCapacity:   config.InitialCapacity,
					StorageSku:        filesystems.StorageSku(config.StorageSku),
					UserDetails: filesystems.UserDetails{
						Email: config.UserEmailAddress,
					},
					MarketplaceDetails: filesystems.MarketplaceDetails{
						OfferId:     offerId,
						PlanId:      config.MarketplacePlanId,
						PublisherId: publisherId,
					},
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r FileSystemResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Qumulo.FileSystemsClient

			id, err := filesystems.ParseFileSystemID(metadata.ResourceData.Id())
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

			var config FileSystemResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			config.Name = id.FileSystemName
			config.ResourceGroupName = id.ResourceGroupName

			if model := resp.Model; model != nil {
				// model.Properties and model.Properties.MarketplaceDetails is not pointer, so we don't need to check nil
				config.Zone = pointer.From(model.Properties.AvailabilityZone)
				config.InitialCapacity = model.Properties.InitialCapacity
				config.MarketplacePlanId = model.Properties.MarketplaceDetails.PlanId
				config.StorageSku = string(model.Properties.StorageSku)
				config.Location = location.Normalize(model.Location)
				config.Tags = tags.Flatten(model.Tags)

				subnetId, err := commonids.ParseSubnetIDInsensitively(model.Properties.DelegatedSubnetId)
				if err != nil {
					return err
				}
				config.SubnetId = subnetId.ID()
			}

			return metadata.Encode(&config)
		},
	}
}

func (r FileSystemResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Qumulo.FileSystemsClient

			id, err := filesystems.ParseFileSystemID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r FileSystemResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Qumulo.FileSystemsClient

			id, err := filesystems.ParseFileSystemID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config FileSystemResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			payload := filesystems.FileSystemResourceUpdate{}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = tags.Expand(config.Tags)
			}

			if _, err := client.Update(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func checkSubnet(ctx context.Context, rawSubnetId string, metadata sdk.ResourceMetaData) error {
	const (
		delegationAction = "Microsoft.Network/virtualNetworks/subnets/join/action"
		delegationName   = "Qumulo.Storage/fileSystems"
	)

	subnetId, err := commonids.ParseSubnetID(rawSubnetId)
	if err != nil {
		return err
	}

	subnet, err := metadata.Client.Network.Subnets.Get(ctx, *subnetId, subnets.GetOperationOptions{})
	if err != nil {
		return fmt.Errorf("checking the subnet: %+v", err)
	}

	if subnet.Model != nil && subnet.Model.Properties != nil && subnet.Model.Properties.Delegations != nil {
		for _, delegation := range *subnet.Model.Properties.Delegations {
			if delegation.Properties != nil && delegation.Properties.Actions != nil &&
				delegation.Properties.ServiceName != nil && strings.EqualFold(*delegation.Properties.ServiceName, delegationName) {
				for _, action := range *delegation.Properties.Actions {
					if strings.EqualFold(action, delegationAction) {
						return nil
					}
				}
			}
		}
	}

	return fmt.Errorf("subnet %q is not delegated %q to %q", rawSubnetId, delegationAction, delegationName)
}
