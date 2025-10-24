package iotoperations

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/instance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type InstanceResource struct{}

var _ sdk.ResourceWithUpdate = InstanceResource{}

type InstanceModel struct {
	Name                 string            `tfschema:"name"`
	ResourceGroupName    string            `tfschema:"resource_group_name"`
	Location             string            `tfschema:"location"`
	Description          *string           `tfschema:"description"`
	Version              *string           `tfschema:"version"`
	ProvisioningState    *string           `tfschema:"provisioning_state"`
	ExtendedLocationName *string           `tfschema:"extended_location_name"`
	ExtendedLocationType *string           `tfschema:"extended_location_type"`
	SchemaRegistryRef    string            `tfschema:"schema_registry_ref"`
	Tags                 map[string]string `tfschema:"tags"`
}

func (r InstanceResource) ModelObject() interface{} {
	return &InstanceModel{}
}

func (r InstanceResource) ResourceType() string {
	return "azurerm_iotoperations_instance"
}

func (r InstanceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return instance.ValidateInstanceID
}

func (r InstanceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 63),
				validation.StringMatch(regexp.MustCompile("^[a-z0-9][a-z0-9-]*[a-z0-9]$"), "Must match ^[a-z0-9][a-z0-9-]*[a-z0-9]$"),
			),
		},
		"resource_group_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 90),
		},
		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"schema_registry_ref": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},
		"extended_location_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"extended_location_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},
	}
}

func (r InstanceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"provisioning_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r InstanceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTOperations.InstanceClient

			var model InstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			id := instance.NewInstanceID(subscriptionId, model.ResourceGroupName, model.Name)

			// Build extended location if provided
			var extendedLocation *instance.ExtendedLocation
			if model.ExtendedLocationName != nil && model.ExtendedLocationType != nil {
				extendedLocation = &instance.ExtendedLocation{
					Name: *model.ExtendedLocationName,
					Type: instance.ExtendedLocationType(*model.ExtendedLocationType),
				}
			}

			// Build properties
			props := &instance.InstanceProperties{
				Description: model.Description,
				Version:     model.Version,
				SchemaRegistryRef: instance.SchemaRegistryRef{
					ResourceId: model.SchemaRegistryRef,
				},
			}

			payload := instance.InstanceResource{
				Location:         model.Location,
				ExtendedLocation: *extendedLocation, // Required field
				Properties:       props,
				Tags:             &model.Tags,
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r InstanceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTOperations.InstanceClient

			id, err := instance.ParseInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			model := InstanceModel{
				Name:              id.InstanceName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if resp.Model != nil {
				model.Location = resp.Model.Location

				if resp.Model.ExtendedLocation.Name != "" {
					model.ExtendedLocationName = &resp.Model.ExtendedLocation.Name
					model.ExtendedLocationType = (*string)(&resp.Model.ExtendedLocation.Type)
				}

				if resp.Model.Properties != nil {
					model.Description = resp.Model.Properties.Description
					model.Version = resp.Model.Properties.Version
					model.SchemaRegistryRef = resp.Model.Properties.SchemaRegistryRef.ResourceId

					if resp.Model.Properties.ProvisioningState != nil {
						provState := string(*resp.Model.Properties.ProvisioningState)
						model.ProvisioningState = &provState
					}
				}

				if resp.Model.Tags != nil {
					model.Tags = *resp.Model.Tags
				}
			}

			return metadata.Encode(&model)
		},
	}
}

func (r InstanceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTOperations.InstanceClient

			id, err := instance.ParseInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model InstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// Check if anything actually changed
			if !metadata.ResourceData.HasChange("tags") {
				return nil
			}

			payload := instance.InstancePatchModel{}
			hasChanges := false

			// Only tags can be updated via PATCH
			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = &model.Tags
				hasChanges = true
			}

			if !hasChanges {
				return nil
			}

			if _, err := client.Update(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r InstanceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTOperations.InstanceClient

			id, err := instance.ParseInstanceID(metadata.ResourceData.Id())
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

// Helper to get pointer to any type
func toPtr[T any](v T) *T {
	return &v
}

// Using your existing toPtr helper (create reverse function)
func fromPtr[T any](ptr *T) T {
	if ptr == nil {
		var zero T
		return zero
	}
	return *ptr
}
