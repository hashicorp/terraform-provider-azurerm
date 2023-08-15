package qumulo

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/qumulostorage/2022-10-12/filesystems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/qumulo/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = FileSystemResource{}
var _ sdk.ResourceWithUpdate = FileSystemResource{}

const (
	offerId     = "qumulo-saas-mpp"
	planId      = "qumulo-on-azure-v1%%gmz7xq9ge3py%%P1M"
	publisherId = "qumulo1584033880660"
)

type FileSystemResource struct{}

func (r FileSystemResource) ModelObject() interface{} {
	return &FileSystemResourceSchema{}
}

type FileSystemResourceSchema struct {
	AvailabilityZone  string                 `tfschema:"availability_zone"`
	AdminPassword     string                 `tfschema:"admin_password"`
	DelegatedSubnetId string                 `tfschema:"delegated_subnet_id"`
	InitialCapacity   int64                  `tfschema:"initial_capacity"`
	Location          string                 `tfschema:"location"`
	Name              string                 `tfschema:"name"`
	ResourceGroupName string                 `tfschema:"resource_group_name"`
	UserEmailAddress  string                 `tfschema:"user_email_address"`
	StorageSku        string                 `tfschema:"storage_sku"`
	Tags              map[string]interface{} `tfschema:"tags"`
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

		"delegated_subnet_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"initial_capacity": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(18, 1000),
		},

		"storage_sku": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(filesystems.PossibleValuesForStorageSku(), false),
		},

		"user_email_address": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"availability_zone": commonschema.ZoneSingleOptionalForceNew(),

		"location": commonschema.Location(),

		"tags": commonschema.Tags(),
	}
}

func (r FileSystemResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r FileSystemResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
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

			payload := filesystems.FileSystemResource{
				Location: location.Normalize(config.Location),
				Tags:     tags.Expand(config.Tags),
				Properties: filesystems.FileSystemResourceProperties{
					AdminPassword:     config.AdminPassword,
					AvailabilityZone:  pointer.To(config.AvailabilityZone),
					DelegatedSubnetId: config.DelegatedSubnetId,
					InitialCapacity:   config.InitialCapacity,
					StorageSku:        filesystems.StorageSku(config.StorageSku),
					UserDetails: filesystems.UserDetails{
						Email: config.UserEmailAddress,
					},
					MarketplaceDetails: filesystems.MarketplaceDetails{
						OfferId:     offerId,
						PlanId:      planId,
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
			schema := FileSystemResourceSchema{}

			id, err := filesystems.ParseFileSystemID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config FileSystemResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			schema.AdminPassword = config.AdminPassword
			schema.UserEmailAddress = config.UserEmailAddress

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if model := resp.Model; model != nil {
				schema.Name = id.FileSystemName
				schema.ResourceGroupName = id.ResourceGroupName
				schema.AvailabilityZone = pointer.From(model.Properties.AvailabilityZone)
				schema.InitialCapacity = model.Properties.InitialCapacity
				schema.StorageSku = string(model.Properties.StorageSku)
				schema.Location = location.Normalize(model.Location)
				schema.Tags = tags.Flatten(model.Tags)

				subnetId, err := commonids.ParseSubnetID(model.Properties.DelegatedSubnetId)
				if err != nil {
					return err
				}
				schema.DelegatedSubnetId = subnetId.ID()
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r FileSystemResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
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
