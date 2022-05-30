package fluidrelay

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/fluidrelay/sdk/2022-04-21/fluidrelayservers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/fluidrelay/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/msi/sdk/2018-11-30/managedidentity"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

//type Identity struct {
//	Type        string `tfschema:"type"`
//	PrincipalID string `tfschema:"principal_id"`
//	TenantID    string `tfschea:"tenant_id"`
//}

type UserAssignedIdentity struct {
	IdentityID  string `tfschema:"identity_id"`
	ClientID    string `tfschema:"client_id"`
	PrincipalID string `tfschema:"principal_id"`
}

type Encryption struct {
	IdentityType        string `tfschema:"identity_type"`
	IdentityResourceId  string `tfschema:"identity_resource_id"`
	KeyEncryptionKeyUrl string `tfschema:"key_encryption_key_url"`
}

type ServersModel struct {
	Name                   string                 `tfschema:"name"`
	ResourceGroup          string                 `tfschema:"resource_group_name"`
	Location               string                 `tfschema:"location"`
	Tags                   map[string]string      `tfschema:"tags"`
	FrsTenantId            string                 `tfschema:"frs_tenant_id"`
	ProvisioningState      string                 `tfschema:"provisioning_state"`
	OrdererEndpoints       []string               `tfschema:"orderer_endpoints"`
	StorageEndpoints       []string               `tfschema:"storage_endpoints"`
	IdentityType           string                 `tfschema:"identity_type"`
	TenantID               string                 `tfschema:"tenant_id"`
	PrincipalID            string                 `tfschema:"principal_id"`
	UserAssignedIdentities []UserAssignedIdentity `tfschema:"user_assigned_identity"`
	Encryption             []Encryption           `tfschema:"encryption"`
}

func (s *ServersModel) GenUserIdentities() *identity.SystemAndUserAssignedMap {
	res := &identity.SystemAndUserAssignedMap{
		Type: "None",
	}
	if s == nil || len(s.UserAssignedIdentities) == 0 {
		return res
	}
	if s.IdentityType != "" {
		res.Type = identity.Type(s.IdentityType)
	}
	if s.PrincipalID != "" {
		res.PrincipalId = s.PrincipalID
	}
	if s.TenantID != "" {
		res.TenantId = s.TenantID
	}
	res.IdentityIds = map[string]identity.UserAssignedIdentityDetails{}
	for _, id := range s.UserAssignedIdentities {
		res.IdentityIds[id.IdentityID] = identity.UserAssignedIdentityDetails{
			ClientId:    utils.Ptr(id.ClientID),
			PrincipalId: utils.Ptr(id.PrincipalID),
		}
	}
	return res
}

func (s *ServersModel) GenEncryption() *fluidrelayservers.EncryptionProperties {
	if len(s.Encryption) == 0 {
		return nil
	}
	encryption := s.Encryption[0]
	res := &fluidrelayservers.EncryptionProperties{
		CustomerManagedKeyEncryption: &fluidrelayservers.CustomerManagedKeyEncryptionProperties{
			KeyEncryptionKeyIdentity: &fluidrelayservers.CustomerManagedKeyEncryptionPropertiesKeyEncryptionKeyIdentity{
				IdentityType:                   utils.Ptr(fluidrelayservers.CmkIdentityType(encryption.IdentityType)),
				UserAssignedIdentityResourceId: utils.Ptr(encryption.IdentityResourceId),
			},
			KeyEncryptionKeyUrl: utils.Ptr(encryption.KeyEncryptionKeyUrl),
		}}
	return res
}

type Servers struct{}

var _ sdk.ResourceWithUpdate = (*Servers)(nil)

func (s Servers) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			Description:  "The Fluid Relay server resource name",
			ValidateFunc: validate.FluidRelayServerName,
		},
		"resource_group_name": commonschema.ResourceGroupName(),
		"location":            commonschema.Location(),
		"tags":                commonschema.Tags(),
		"identity_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			//Default:  "SystemAssigned",
			ValidateFunc: validation.StringInSlice([]string{
				"SystemAssigned",
				"UserAssigned",
				"SystemAssigned, UserAssigned",
				"None",
			}, false),
		},
		"user_assigned_identity": {
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Computed:    true,
			ConfigMode:  pluginsdk.SchemaConfigModeBlock,
			Description: "The list of user identities associated with the resource.",
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"identity_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: managedidentity.ValidateUserAssignedIdentitiesID,
					},
					"client_id": {
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Optional:    true,
						Description: "The client id of user assigned identity.",
					},
					"principal_id": {
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Computed:    true,
						Description: "The principal id of user assigned identity.",
					},
				},
			},
		},
		"encryption": {
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Description: "Create with Cmk.",
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"identity_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice(
							fluidrelayservers.PossibleValuesForCmkIdentityType(),
							false),
						Description: "Values can be SystemAssigned or UserAssigned.",
					},
					"key_encryption_key_url": {
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Description: "user assigned identity to use for accessing key encryption key Url. Ex: /subscriptions/fa5fc227-a624-475e-b696-cdd604c735bc/resourceGroups/<resource group>/providers/Microsoft.ManagedIdentity/userAssignedIdentities/myId. Mutually exclusive with identityType systemAssignedIdentity.",
					},
					"identity_resource_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: managedidentity.ValidateUserAssignedIdentitiesID,
					},
				},
			},
		},
	}
}

func (s Servers) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"frs_tenant_id": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The Fluid tenantId for this server",
		},
		"provisioning_state": {
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
		"principal_id": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The principal ID of resource identity",
		},
		"tenant_id": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The tenant ID of resource",
		},
	}
}

func (s Servers) ModelObject() interface{} {
	return &ServersModel{}
}

func (s Servers) ResourceType() string {
	return "azurerm_fluid_relay_servers"
}

func (s Servers) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) (err error) {
			client := meta.Client.FluidRelay.ServerClient

			var model ServersModel
			if err = meta.Decode(&model); err != nil {
				return err
			}

			account := meta.Client.Account
			id := fluidrelayservers.NewFluidRelayServerID(account.SubscriptionId, model.ResourceGroup, model.Name)

			existing, err := client.Get(ctx, id)
			if !response.WasNotFound(existing.HttpResponse) {
				if err != nil {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				return meta.ResourceRequiresImport(s.ResourceType(), id)
			}

			serverReq := fluidrelayservers.FluidRelayServer{
				Location: azure.NormalizeLocation(model.Location),
				Name:     utils.Ptr(model.Name),
			}
			serverReq.Tags = utils.TryPtr(model.Tags)
			serverReq.Properties = &fluidrelayservers.FluidRelayServerProperties{}
			serverReq.Properties.Encryption = model.GenEncryption()
			serverReq.Identity = model.GenUserIdentities()

			log := fmt.Sprintf("generate req: %s by model: %s", utils.JSONStr(serverReq), utils.JSONStr(model))
			meta.Logger.Infof("start creating: %s, data: %s", id, log)
			//return fmt.Errorf("interrupt by test: %s", log)
			resp, err := client.CreateOrUpdate(ctx, id, serverReq)
			if err != nil {
				return fmt.Errorf("creating %v err: %+v, resp: %+v", id, err, resp)
			}
			meta.SetID(id)

			return nil
		},
	}
}

// Update tags
func (s Servers) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) (err error) {
			client := meta.Client.FluidRelay.ServerClient
			id, err := fluidrelayservers.ParseFluidRelayServerIDInsensitively(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ServersModel
			if err = meta.Decode(&model); err != nil {
				return fmt.Errorf("decoding err: %+v", err)
			}

			var upd fluidrelayservers.FluidRelayServerUpdate
			if meta.ResourceData.HasChanges("tags", "location") {
				upd.Tags = &model.Tags
				upd.Location = utils.String(azure.NormalizeLocation(model.Location))
			}
			if meta.ResourceData.HasChanges("identity_type", "user_assigned_identity") {
				upd.Identity = model.GenUserIdentities()
			}
			if meta.ResourceData.HasChanges("encryption") {
				upd.Properties = &fluidrelayservers.FluidRelayServerUpdateProperties{}
				upd.Properties.Encryption = model.GenEncryption()
			}
			_, err = client.Update(ctx, *id, upd)
			if err != nil {
				return fmt.Errorf("updating fluid relay err: %v", err)
			}
			// do a read after update
			return nil
		},
	}
}

func (s Servers) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.FluidRelay.ServerClient

			id, err := fluidrelayservers.ParseFluidRelayServerID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			if id == nil {
				return fmt.Errorf("parsed id is nil")
			}

			server, err := client.Get(ctx, *id)
			if err != nil {
				return err
			}

			model := server.Model
			if model == nil {
				return fmt.Errorf("got fluid relay server as nil")
			}

			output := ServersModel{
				Name:          id.FluidRelayServerName,
				ResourceGroup: id.ResourceGroup,
				Location:      model.Location,
				IdentityType:  string(model.Identity.Type),
			}
			if model.Identity != nil {
				output.TenantID = model.Identity.TenantId
				output.PrincipalID = model.Identity.PrincipalId
				for id, details := range model.Identity.IdentityIds {
					// try parse id, because the response id could be modified to lower-case
					iid, err := commonids.ParseUserAssignedIdentityIDInsensitively(id)
					if err != nil {
						meta.Logger.Warnf("normalize managed identity id `%s` parse err: %v", id, err)
					} else {
						id = iid.ID()
					}
					output.UserAssignedIdentities = append(output.UserAssignedIdentities, UserAssignedIdentity{
						IdentityID:  id,
						ClientID:    utils.Value(details.ClientId),
						PrincipalID: utils.Value(details.PrincipalId),
					})
				}
			}
			output.Tags = utils.Value(server.Model.Tags)
			if prop := model.Properties; prop != nil {
				if prop.FrsTenantId != nil {
					output.FrsTenantId = *prop.FrsTenantId
				}
				if prop.ProvisioningState != nil {
					output.ProvisioningState = string(*prop.ProvisioningState)
				}
				if points := prop.FluidRelayEndpoints; points != nil {
					if points.OrdererEndpoints != nil {
						output.OrdererEndpoints = *points.OrdererEndpoints
					}
					if points.StorageEndpoints != nil {
						output.StorageEndpoints = *points.StorageEndpoints
					}
				}
			}
			return meta.Encode(&output)
		},
	}
}

func (s Servers) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.FluidRelay.ServerClient

			id, err := fluidrelayservers.ParseFluidRelayServerID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			meta.Logger.Infof("deleting %s", id)
			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("delete %s err: %v", id, err)
			}
			return nil
		},
	}
}

func (s Servers) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return fluidrelayservers.ValidateFluidRelayServerID
}
