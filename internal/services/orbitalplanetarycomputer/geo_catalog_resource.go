// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package orbitalplanetarycomputer

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/orbitalplanetarycomputer/2026-04-15/geocatalogs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name geo_catalog -service-package-name orbitalplanetarycomputer -properties "name,resource_group_name" -known-values "subscription_id:data.Subscriptions.Primary"

type GeoCatalogResource struct{}

var (
	_ sdk.ResourceWithUpdate        = GeoCatalogResource{}
	_ sdk.ResourceWithIdentity      = GeoCatalogResource{}
	_ sdk.ResourceWithCustomizeDiff = GeoCatalogResource{}
)

func (r GeoCatalogResource) Identity() resourceids.ResourceId {
	return &geocatalogs.GeoCatalogId{}
}

type GeoCatalogModel struct {
	Name              string                                     `tfschema:"name"`
	ResourceGroupName string                                     `tfschema:"resource_group_name"`
	Location          string                                     `tfschema:"location"`
	Identity          []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Tags              map[string]string                          `tfschema:"tags"`
}

func (r GeoCatalogResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z0-9-]{3,24}$"),
				"GeoCatalog name must consist only of letters, numbers, hyphens (-), and have length between 3 and 24 characters (inclusive)",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"tags": commonschema.Tags(),
	}
}

func (r GeoCatalogResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r GeoCatalogResource) ModelObject() interface{} {
	return &GeoCatalogModel{}
}

func (r GeoCatalogResource) ResourceType() string {
	return "azurerm_geo_catalog"
}

func (r GeoCatalogResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.OrbitalPlanetaryComputer.GeoCatalogsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model GeoCatalogModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := geocatalogs.NewGeoCatalogID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			param := geocatalogs.GeoCatalog{
				Location: location.Normalize(model.Location),
				Identity: expandedIdentity,
				Tags:     pointer.To(model.Tags),
				Properties: &geocatalogs.GeoCatalogProperties{
					Tier: pointer.To(geocatalogs.CatalogTierBasic),
				},
			}

			if err := client.CreateCallbackThenPoll(ctx, id, param, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			return nil
		},
	}
}

func (r GeoCatalogResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.OrbitalPlanetaryComputer.GeoCatalogsClient

			id, err := geocatalogs.ParseGeoCatalogID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model GeoCatalogModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			param := geocatalogs.GeoCatalogUpdate{}

			if metadata.ResourceData.HasChange("identity") {
				param.Identity = r.expandIdentity(model.Identity)
			}

			if metadata.ResourceData.HasChange("tags") {
				param.Tags = pointer.To(model.Tags)
			}

			if err := client.UpdateThenPoll(ctx, *id, param); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r GeoCatalogResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.OrbitalPlanetaryComputer.GeoCatalogsClient

			id, err := geocatalogs.ParseGeoCatalogID(metadata.ResourceData.Id())
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

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (r GeoCatalogResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.OrbitalPlanetaryComputer.GeoCatalogsClient

			id, err := geocatalogs.ParseGeoCatalogID(metadata.ResourceData.Id())
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

func (r GeoCatalogResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return geocatalogs.ValidateGeoCatalogID
}

func (r GeoCatalogResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if rawIdentityIds, ok := metadata.ResourceDiff.GetOk("identity.0.identity_ids"); ok {
				identityIds := rawIdentityIds.(*pluginsdk.Set)
				identityType := metadata.ResourceDiff.Get("identity.0.type").(string)
				if identityType == string(identity.TypeSystemAssignedUserAssigned) {
					return fmt.Errorf("%q resource identity is not supported", string(identity.TypeSystemAssignedUserAssigned))
				} else if identityIds.Len() > 0 && identityType != string(identity.TypeUserAssigned) {
					return fmt.Errorf("`identity_ids` can only be specified when `type` is set to %q", string(identity.TypeUserAssigned))
				}
			}

			return nil
		},
	}
}

func (r GeoCatalogResource) expandIdentity(input []identity.ModelSystemAssignedUserAssigned) *geocatalogs.ManagedServiceIdentityUpdate {
	managedServiceIdentityUpdate := &geocatalogs.ManagedServiceIdentityUpdate{
		Type:                   pointer.To(geocatalogs.ManagedServiceIdentityTypeNone),
		UserAssignedIdentities: &map[string]geocatalogs.UserAssignedIdentity{},
	}

	if len(input) == 0 {
		return managedServiceIdentityUpdate
	}

	managedServiceIdentityUpdate.Type = pointer.ToEnum[geocatalogs.ManagedServiceIdentityType](string(input[0].Type))

	for _, identityId := range input[0].IdentityIds {
		(*managedServiceIdentityUpdate.UserAssignedIdentities)[identityId] = geocatalogs.UserAssignedIdentity{}
	}

	return managedServiceIdentityUpdate
}

func (r GeoCatalogResource) flatten(metadata sdk.ResourceMetaData, id *geocatalogs.GeoCatalogId, model *geocatalogs.GeoCatalog) error {
	state := GeoCatalogModel{
		Name:              id.GeoCatalogName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model != nil {
		state.Location = location.Normalize(model.Location)

		flattenedIdentity, err := identity.FlattenLegacySystemAndUserAssignedMapToModel(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		state.Identity = flattenedIdentity

		state.Tags = pointer.From(model.Tags)
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}
