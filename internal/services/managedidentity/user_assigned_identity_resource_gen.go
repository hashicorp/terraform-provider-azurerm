package managedidentity

// NOTE: this file is generated - manual changes will be overwritten.
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.
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
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2023-01-31/managedidentities"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = UserAssignedIdentityResource{}
var _ sdk.ResourceWithUpdate = UserAssignedIdentityResource{}

type UserAssignedIdentityResource struct{}

func (r UserAssignedIdentityResource) ModelObject() interface{} {
	return &UserAssignedIdentityResourceSchema{}
}

type UserAssignedIdentityResourceSchema struct {
	ClientId          string                 `tfschema:"client_id"`
	Location          string                 `tfschema:"location"`
	Name              string                 `tfschema:"name"`
	PrincipalId       string                 `tfschema:"principal_id"`
	ResourceGroupName string                 `tfschema:"resource_group_name"`
	Tags              map[string]interface{} `tfschema:"tags"`
	TenantId          string                 `tfschema:"tenant_id"`
}

func (r UserAssignedIdentityResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateUserAssignedIdentityID
}
func (r UserAssignedIdentityResource) ResourceType() string {
	return "azurerm_user_assigned_identity"
}
func (r UserAssignedIdentityResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),
		"name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"resource_group_name": commonschema.ResourceGroupName(),
		"tags":                commonschema.Tags(),
	}
}
func (r UserAssignedIdentityResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"client_id": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"principal_id": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"tenant_id": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
	}
}
func (r UserAssignedIdentityResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedIdentity.V20230131.ManagedIdentities

			var config UserAssignedIdentityResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId

			id := commonids.NewUserAssignedIdentityID(subscriptionId, config.ResourceGroupName, config.Name)

			existing, err := client.UserAssignedIdentitiesGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			var payload managedidentities.Identity
			if err := r.mapUserAssignedIdentityResourceSchemaToIdentity(config, &payload); err != nil {
				return fmt.Errorf("mapping schema model to sdk model: %+v", err)
			}

			if _, err := client.UserAssignedIdentitiesCreateOrUpdate(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
func (r UserAssignedIdentityResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedIdentity.V20230131.ManagedIdentities
			schema := UserAssignedIdentityResourceSchema{}

			id, err := commonids.ParseUserAssignedIdentityID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.UserAssignedIdentitiesGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				schema.Name = id.UserAssignedIdentityName
				schema.ResourceGroupName = id.ResourceGroupName
				if err := r.mapIdentityToUserAssignedIdentityResourceSchema(*model, &schema); err != nil {
					return fmt.Errorf("flattening model: %+v", err)
				}
			}

			return metadata.Encode(&schema)
		},
	}
}
func (r UserAssignedIdentityResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedIdentity.V20230131.ManagedIdentities

			id, err := commonids.ParseUserAssignedIdentityID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.UserAssignedIdentitiesDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
func (r UserAssignedIdentityResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedIdentity.V20230131.ManagedIdentities

			id, err := commonids.ParseUserAssignedIdentityID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config UserAssignedIdentityResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			var payload managedidentities.IdentityUpdate
			if err := r.mapUserAssignedIdentityResourceSchemaToIdentityUpdate(config, &payload); err != nil {
				return fmt.Errorf("mapping schema model to sdk model: %+v", err)
			}

			if _, err := client.UserAssignedIdentitiesUpdate(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r UserAssignedIdentityResource) mapUserAssignedIdentityResourceSchemaToUserAssignedIdentityProperties(input UserAssignedIdentityResourceSchema, output *managedidentities.UserAssignedIdentityProperties) error {

	return nil
}

func (r UserAssignedIdentityResource) mapUserAssignedIdentityPropertiesToUserAssignedIdentityResourceSchema(input managedidentities.UserAssignedIdentityProperties, output *UserAssignedIdentityResourceSchema) error {
	output.ClientId = pointer.From(input.ClientId)
	output.PrincipalId = pointer.From(input.PrincipalId)
	output.TenantId = pointer.From(input.TenantId)
	return nil
}

func (r UserAssignedIdentityResource) mapUserAssignedIdentityResourceSchemaToIdentity(input UserAssignedIdentityResourceSchema, output *managedidentities.Identity) error {
	output.Location = location.Normalize(input.Location)
	output.Tags = tags.Expand(input.Tags)

	if output.Properties == nil {
		output.Properties = &managedidentities.UserAssignedIdentityProperties{}
	}
	if err := r.mapUserAssignedIdentityResourceSchemaToUserAssignedIdentityProperties(input, output.Properties); err != nil {
		return fmt.Errorf("mapping Schema to SDK Field %q / Model %q: %+v", "UserAssignedIdentityProperties", "Properties", err)
	}

	return nil
}

func (r UserAssignedIdentityResource) mapIdentityToUserAssignedIdentityResourceSchema(input managedidentities.Identity, output *UserAssignedIdentityResourceSchema) error {
	output.Location = location.Normalize(input.Location)
	output.Tags = tags.Flatten(input.Tags)

	if input.Properties == nil {
		input.Properties = &managedidentities.UserAssignedIdentityProperties{}
	}
	if err := r.mapUserAssignedIdentityPropertiesToUserAssignedIdentityResourceSchema(*input.Properties, output); err != nil {
		return fmt.Errorf("mapping SDK Field %q / Model %q to Schema: %+v", "UserAssignedIdentityProperties", "Properties", err)
	}

	return nil
}

func (r UserAssignedIdentityResource) mapUserAssignedIdentityResourceSchemaToIdentityUpdate(input UserAssignedIdentityResourceSchema, output *managedidentities.IdentityUpdate) error {
	output.Tags = tags.Expand(input.Tags)

	if output.Properties == nil {
		output.Properties = &managedidentities.UserAssignedIdentityProperties{}
	}
	if err := r.mapUserAssignedIdentityResourceSchemaToUserAssignedIdentityProperties(input, output.Properties); err != nil {
		return fmt.Errorf("mapping Schema to SDK Field %q / Model %q: %+v", "UserAssignedIdentityProperties", "Properties", err)
	}

	return nil
}

func (r UserAssignedIdentityResource) mapIdentityUpdateToUserAssignedIdentityResourceSchema(input managedidentities.IdentityUpdate, output *UserAssignedIdentityResourceSchema) error {
	output.Tags = tags.Flatten(input.Tags)

	if input.Properties == nil {
		input.Properties = &managedidentities.UserAssignedIdentityProperties{}
	}
	if err := r.mapUserAssignedIdentityPropertiesToUserAssignedIdentityResourceSchema(*input.Properties, output); err != nil {
		return fmt.Errorf("mapping SDK Field %q / Model %q to Schema: %+v", "UserAssignedIdentityProperties", "Properties", err)
	}

	return nil
}
