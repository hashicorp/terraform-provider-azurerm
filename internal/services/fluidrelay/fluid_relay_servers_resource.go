// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fluidrelay

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/fluidrelay/2022-05-26/fluidrelayservers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/fluidrelay/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ServerModel struct {
	Name             string                                     `tfschema:"name"`
	ResourceGroup    string                                     `tfschema:"resource_group_name"`
	Location         string                                     `tfschema:"location"`
	StorageSKU       string                                     `tfschema:"storage_sku"`
	FrsTenantId      string                                     `tfschema:"frs_tenant_id"`
	PrimaryKey       string                                     `tfschema:"primary_key"`
	SecondaryKey     string                                     `tfschema:"secondary_key"`
	OrdererEndpoints []string                                   `tfschema:"orderer_endpoints"`
	StorageEndpoints []string                                   `tfschema:"storage_endpoints"`
	ServiceEndpoints []string                                   `tfschema:"service_endpoints"`
	Tags             map[string]string                          `tfschema:"tags"`
	Identity         []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
}

func (s *ServerModel) flattenIdentity(input *identity.SystemAndUserAssignedMap) error {
	if input == nil {
		return nil
	}
	config := identity.SystemAndUserAssignedMap{
		Type:        input.Type,
		PrincipalId: input.PrincipalId,
		TenantId:    input.TenantId,
		IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
	}
	for k, v := range input.IdentityIds {
		config.IdentityIds[k] = identity.UserAssignedIdentityDetails{
			ClientId:    v.ClientId,
			PrincipalId: v.PrincipalId,
		}
	}
	model, err := identity.FlattenSystemAndUserAssignedMapToModel(&config)
	if err != nil {
		return err
	}
	s.Identity = *model
	return nil
}

type Server struct{}

var _ sdk.ResourceWithUpdate = (*Server)(nil)

func (s Server) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FluidRelayServerName,
		},
		"resource_group_name": commonschema.ResourceGroupName(),
		"location":            commonschema.Location(),
		"tags":                commonschema.Tags(),
		"identity":            commonschema.SystemAssignedUserAssignedIdentityOptional(),
		"storage_sku": {
			// todo remove computed when https://github.com/Azure/azure-rest-api-specs/issues/19700 is fixed
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(fluidrelayservers.PossibleValuesForStorageSKU(), false),
		},
	}
}

func (s Server) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"frs_tenant_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"orderer_endpoints": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"storage_endpoints": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"service_endpoints": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"primary_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},
	}
}

func (s Server) ModelObject() interface{} {
	return &ServerModel{}
}

func (s Server) ResourceType() string {
	return "azurerm_fluid_relay_server"
}

func (s Server) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) (err error) {
			client := meta.Client.FluidRelay.FluidRelayServers

			var model ServerModel
			if err = meta.Decode(&model); err != nil {
				return err
			}

			id := fluidrelayservers.NewFluidRelayServerID(meta.Client.Account.SubscriptionId, model.ResourceGroup, model.Name)

			existing, err := client.Get(ctx, id)
			if !response.WasNotFound(existing.HttpResponse) {
				if err != nil {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				return meta.ResourceRequiresImport(s.ResourceType(), id)
			}

			payload := fluidrelayservers.FluidRelayServer{
				Location:   location.Normalize(model.Location),
				Properties: &fluidrelayservers.FluidRelayServerProperties{},
				Tags:       &model.Tags,
			}
			payload.Identity, err = identity.ExpandSystemAndUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding user identities: %+v", err)
			}

			if model.StorageSKU != "" {
				payload.Properties.Storagesku = pointer.To(fluidrelayservers.StorageSKU(model.StorageSKU))
			}

			if _, err = client.CreateOrUpdate(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %v err: %+v", id, err)
			}
			meta.SetID(id)

			return nil
		},
	}
}

func (s Server) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) (err error) {
			client := meta.Client.FluidRelay.FluidRelayServers
			id, err := fluidrelayservers.ParseFluidRelayServerID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ServerModel
			if err = meta.Decode(&model); err != nil {
				return fmt.Errorf("decoding err: %+v", err)
			}

			payload := fluidrelayservers.FluidRelayServerUpdate{}
			if meta.ResourceData.HasChange("tags") {
				payload.Tags = &model.Tags
			}
			if meta.ResourceData.HasChange("identity") {
				payload.Identity, err = identity.ExpandSystemAndUserAssignedMapFromModel(model.Identity)
				if err != nil {
					return fmt.Errorf("expanding user identities: %+v", err)
				}
			}
			if _, err = client.Update(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %v", id, err)
			}

			return nil
		},
	}
}

func (s Server) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.FluidRelay.FluidRelayServers

			id, err := fluidrelayservers.ParseFluidRelayServerID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			server, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(server.HttpResponse) {
					return meta.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			keys, err := client.ListKeys(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving keys for %s: %+v", *id, err)
			}

			output := &ServerModel{
				Name:          id.FluidRelayServerName,
				ResourceGroup: id.ResourceGroup,
			}
			if model := server.Model; model != nil {
				output.Location = location.Normalize(model.Location)

				if err = output.flattenIdentity(model.Identity); err != nil {
					return fmt.Errorf("flattening `identity`: %v", err)
				}
				if server.Model.Tags != nil {
					output.Tags = *server.Model.Tags
				}
				if prop := model.Properties; prop != nil {
					if prop.FrsTenantId != nil {
						output.FrsTenantId = *prop.FrsTenantId
					}
					if points := prop.FluidRelayEndpoints; points != nil {
						if points.OrdererEndpoints != nil {
							output.OrdererEndpoints = *points.OrdererEndpoints
						}
						if points.StorageEndpoints != nil {
							output.StorageEndpoints = *points.StorageEndpoints
						}

						if points.ServiceEndpoints != nil {
							output.ServiceEndpoints = *points.ServiceEndpoints
						}
					}
				}
				if val, ok := meta.ResourceData.GetOk("storage_sku"); ok {
					output.StorageSKU = val.(string)
				}
			}

			if model := keys.Model; model != nil {
				output.PrimaryKey = utils.NormalizeNilableString(model.Key1)
				output.SecondaryKey = utils.NormalizeNilableString(model.Key2)
			}

			return meta.Encode(output)
		},
	}
}

func (s Server) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.FluidRelay.FluidRelayServers

			id, err := fluidrelayservers.ParseFluidRelayServerID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			meta.Logger.Infof("deleting %s", id)
			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %v", id, err)
			}
			return nil
		},
	}
}

func (s Server) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return fluidrelayservers.ValidateFluidRelayServerID
}
