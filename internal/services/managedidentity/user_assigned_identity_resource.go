// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedidentity

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
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2024-11-30/identities"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedidentity/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.Resource                   = UserAssignedIdentityResource{}
	_ sdk.ResourceWithStateMigration = UserAssignedIdentityResource{}
)

type UserAssignedIdentityResource struct{}

func (r UserAssignedIdentityResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.UserAssignedIdentityV0ToV1{},
		},
	}
}

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
			client := metadata.Client.ManagedIdentity.V20241130.Identities

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

			payload := identities.Identity{
				Location:   location.Normalize(config.Location),
				Tags:       tags.Expand(config.Tags),
				Properties: &identities.UserAssignedIdentityProperties{},
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
			client := metadata.Client.ManagedIdentity.V20241130.Identities
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
				schema.Location = location.Normalize(model.Location)
				schema.Tags = tags.Flatten(model.Tags)

				if model.Properties != nil {
					schema.ClientId = pointer.From(model.Properties.ClientId)
					schema.PrincipalId = pointer.From(model.Properties.PrincipalId)
					schema.TenantId = pointer.From(model.Properties.TenantId)
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
			client := metadata.Client.ManagedIdentity.V20241130.Identities

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
			client := metadata.Client.ManagedIdentity.V20241130.Identities

			id, err := commonids.ParseUserAssignedIdentityID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config UserAssignedIdentityResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			payload := identities.IdentityUpdate{
				Tags:       tags.Expand(config.Tags),
				Properties: &identities.UserAssignedIdentityProperties{},
			}

			if _, err := client.UserAssignedIdentitiesUpdate(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
