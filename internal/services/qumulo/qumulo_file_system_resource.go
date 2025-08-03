// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/subnets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/qumulostorage/2024-06-19/filesystems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/qumulo/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = FileSystemResource{}

type FileSystemResource struct{}

func (r FileSystemResource) ModelObject() interface{} {
	return &FileSystemResourceSchema{}
}

type FileSystemResourceSchema struct {
	AdminPassword     string            `tfschema:"admin_password"`
	Location          string            `tfschema:"location"`
	Name              string            `tfschema:"name"`
	OfferId           string            `tfschema:"offer_id"`
	PlanId            string            `tfschema:"plan_id"`
	PublisherId       string            `tfschema:"publisher_id"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	StorageSku        string            `tfschema:"storage_sku"`
	SubnetId          string            `tfschema:"subnet_id"`
	Tags              map[string]string `tfschema:"tags"`
	Email             string            `tfschema:"email"`
	Zone              string            `tfschema:"zone"`
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

		"location": commonschema.Location(),

		"admin_password": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			Sensitive:    true,
			ValidateFunc: validate.ValidatePasswordComplexity,
		},

		"email": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.Email,
		},

		"storage_sku": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			// Use enum in SDK if swagger issue is resolved https://github.com/Azure/azure-rest-api-specs/issues/34017
			ValidateFunc: validation.StringInSlice([]string{
				"Cold_LRS",
				"Hot_LRS",
				"Hot_ZRS",
			}, false),
		},

		"subnet_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"zone": commonschema.ZoneSingleRequiredForceNew(),

		"offer_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      "qumulo-saas-mpp",
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"plan_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      "azure-native-qumulo-v3",
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"publisher_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      "qumulo1584033880660",
			ValidateFunc: validation.StringIsNotEmpty,
		},

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
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config FileSystemResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := filesystems.NewFileSystemID(subscriptionId, config.ResourceGroupName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			// check the subnet is valid: https://learn.microsoft.com/en-us/azure/partner-solutions/qumulo/qumulo-troubleshoot#you-cant-create-a-resource
			if err := checkSubnet(ctx, config.SubnetId, metadata); err != nil {
				return err
			}

			payload := filesystems.LiftrBaseStorageFileSystemResource{
				Location: location.Normalize(config.Location),
				Tags:     pointer.To(config.Tags),
				Properties: &filesystems.LiftrBaseStorageFileSystemResourceProperties{
					AdminPassword:     config.AdminPassword,
					AvailabilityZone:  pointer.To(config.Zone),
					DelegatedSubnetId: config.SubnetId,
					StorageSku:        config.StorageSku,
					UserDetails: filesystems.LiftrBaseUserDetails{
						Email: config.Email,
					},
					MarketplaceDetails: filesystems.LiftrBaseMarketplaceDetails{
						OfferId:     config.OfferId,
						PlanId:      config.PlanId,
						PublisherId: pointer.To(config.PublisherId),
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

			state := FileSystemResourceSchema{
				Name:              id.FileSystemName,
				ResourceGroupName: id.ResourceGroupName,

				AdminPassword: config.AdminPassword,
				Email:         config.Email,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				props := model.Properties
				state.OfferId = props.MarketplaceDetails.OfferId
				state.PlanId = props.MarketplaceDetails.PlanId
				state.PublisherId = pointer.From(props.MarketplaceDetails.PublisherId)
				state.StorageSku = props.StorageSku
				state.Zone = pointer.From(props.AvailabilityZone)

				subnetId, err := commonids.ParseSubnetIDInsensitively(props.DelegatedSubnetId)
				if err != nil {
					return err
				}
				state.SubnetId = subnetId.ID()
			}

			return metadata.Encode(&state)
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

			payload := filesystems.LiftrBaseStorageFileSystemResourceUpdate{}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(config.Tags)
			}

			if _, err := client.Update(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func checkSubnet(ctx context.Context, rawSubnetId string, metadata sdk.ResourceMetaData) error {
	delegationAction := "Microsoft.Network/virtualNetworks/subnets/join/action"
	delegationName := "Qumulo.Storage/fileSystems"

	subnetId, err := commonids.ParseSubnetID(rawSubnetId)
	if err != nil {
		return err
	}

	subnet, err := metadata.Client.Network.Subnets.Get(ctx, *subnetId, subnets.GetOperationOptions{})
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", subnetId, err)
	}

	if subnet.Model != nil && subnet.Model.Properties != nil && subnet.Model.Properties.Delegations != nil {
		for _, delegation := range *subnet.Model.Properties.Delegations {
			if props := delegation.Properties; props != nil && props.Actions != nil &&
				props.ServiceName != nil && strings.EqualFold(*props.ServiceName, delegationName) {
				for _, action := range *props.Actions {
					if strings.EqualFold(action, delegationAction) {
						return nil
					}
				}
			}
		}
	}

	return fmt.Errorf("subnet %q is missing action %q on delegation %q", rawSubnetId, delegationAction, delegationName)
}
