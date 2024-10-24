// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedhsm

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-04-01/roledefinitions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-07-01/managedhsms"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/keyvault/7.4/keyvault"
)

type KeyVaultMHSMRoleDefinitionModel struct {
	ManagedHSMID      string       `tfschema:"managed_hsm_id"`
	Name              string       `tfschema:"name"`
	RoleName          string       `tfschema:"role_name"`
	Description       string       `tfschema:"description"`
	Permission        []Permission `tfschema:"permission"`
	RoleType          string       `tfschema:"role_type"`
	ResourceManagerId string       `tfschema:"resource_manager_id"`

	// TODO: remove in 4.0
	VaultBaseUrl string `tfschema:"vault_base_url,removedInNextMajorVersion"`
}

type Permission struct {
	Actions        []string `tfschema:"actions"`
	NotActions     []string `tfschema:"not_actions"`
	DataActions    []string `tfschema:"data_actions"`
	NotDataActions []string `tfschema:"not_data_actions"`
}

type KeyVaultMHSMRoleDefinitionResource struct{}

var _ sdk.ResourceWithStateMigration = KeyVaultMHSMRoleDefinitionResource{}
var _ sdk.ResourceWithUpdate = KeyVaultMHSMRoleDefinitionResource{}

// Arguments ...
// skip `assignable_scopes` field support as https://github.com/Azure/azure-rest-api-specs/issues/23045
func (r KeyVaultMHSMRoleDefinitionResource) Arguments() map[string]*pluginsdk.Schema {
	s := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"managed_hsm_id": func() *pluginsdk.Schema {
			s := &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ForceNew:     true,
				ValidateFunc: managedhsms.ValidateManagedHSMID,
			}
			if features.FourPointOhBeta() {
				s.Required = true
			} else {
				s.Optional = true
				s.Computed = true
			}
			return s
		}(),

		"role_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"permission": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"actions": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"not_actions": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"data_actions": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(func() (res []string) {
								for _, v := range keyvault.PossibleDataActionValues() {
									res = append(res, string(v))
								}
								return
							}(), false),
						},
						Set: pluginsdk.HashString,
					},

					"not_data_actions": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(func() (res []string) {
								for _, v := range keyvault.PossibleDataActionValues() {
									res = append(res, string(v))
								}
								return
							}(), false),
						},
						Set: pluginsdk.HashString,
					},
				},
			},
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}

	if !features.FourPointOhBeta() {
		s["vault_base_url"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		}
	}

	return s
}

func (r KeyVaultMHSMRoleDefinitionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"role_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"resource_manager_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r KeyVaultMHSMRoleDefinitionResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.ManagedHSMRoleDefinitionV0ToV1{},
		},
	}
}

func (r KeyVaultMHSMRoleDefinitionResource) ModelObject() interface{} {
	return &KeyVaultMHSMRoleDefinitionModel{}
}

func (r KeyVaultMHSMRoleDefinitionResource) ResourceType() string {
	return "azurerm_key_vault_managed_hardware_security_module_role_definition"
}

func (r KeyVaultMHSMRoleDefinitionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedHSMs.DataPlaneRoleDefinitionsClient
			domainSuffix, ok := metadata.Client.Account.Environment.ManagedHSM.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Managed HSM domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			var config KeyVaultMHSMRoleDefinitionModel
			if err := metadata.Decode(&config); err != nil {
				return err
			}

			var managedHsmId *managedhsms.ManagedHSMId
			var endpoint *parse.ManagedHSMDataPlaneEndpoint
			var err error
			if config.ManagedHSMID != "" {
				managedHsmId, err = managedhsms.ParseManagedHSMID(config.ManagedHSMID)
				if err != nil {
					return err
				}
				baseUri, err := metadata.Client.ManagedHSMs.BaseUriForManagedHSM(ctx, *managedHsmId)
				if err != nil {
					return fmt.Errorf("determining the Data Plane Endpoint for %s: %+v", *managedHsmId, err)
				}
				if baseUri == nil {
					return fmt.Errorf("unable to determine the Data Plane Endpoint for %q", *managedHsmId)
				}
				endpoint, err = parse.ManagedHSMEndpoint(*baseUri, domainSuffix)
				if err != nil {
					return fmt.Errorf("parsing the Data Plane Endpoint %q: %+v", *endpoint, err)
				}
			}

			if managedHsmId == nil && !features.FourPointOhBeta() {
				endpoint, err = parse.ManagedHSMEndpoint(config.VaultBaseUrl, domainSuffix)
				if err != nil {
					return fmt.Errorf("parsing the Data Plane Endpoint %q: %+v", *endpoint, err)
				}
				subscriptionId := commonids.NewSubscriptionID(metadata.Client.Account.SubscriptionId)
				managedHsmId, err = metadata.Client.ManagedHSMs.ManagedHSMIDFromBaseUrl(ctx, subscriptionId, endpoint.BaseURI(), domainSuffix)
				if err != nil {
					return fmt.Errorf("determining the Managed HSM ID for %q: %+v", endpoint.BaseURI(), err)
				}
				if managedHsmId == nil {
					return fmt.Errorf("unable to determine the Resource Manager ID")
				}
			}

			// need a lock for hsm subresource create/update/delete, or API may respond error as below
			// Status=409 Code="Conflict" Message="There was a conflict while trying to delete the role assignment.
			locks.ByName(managedHsmId.ID(), "azurerm_key_vault_managed_hardware_security_module")
			defer locks.UnlockByName(managedHsmId.ID(), "azurerm_key_vault_managed_hardware_security_module")

			scope := keyvault.RoleScopeGlobal
			id := parse.NewManagedHSMDataPlaneRoleDefinitionID(endpoint.ManagedHSMName, endpoint.DomainSuffix, string(scope), config.Name)
			existing, err := client.Get(ctx, id.BaseURI(), id.Scope, id.ManagedHSMName)
			if !utils.ResponseWasNotFound(existing.Response) {
				if err != nil {
					return fmt.Errorf("checking for the existence of an existing %q: %+v", id, err)
				}
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := keyvault.RoleDefinitionCreateParameters{
				Properties: &keyvault.RoleDefinitionProperties{
					RoleName:    pointer.To(config.RoleName),
					Description: pointer.To(config.Description),
					RoleType:    keyvault.RoleTypeCustomRole,
					Permissions: expandKeyVaultMHSMRolePermissions(config.Permission),
					AssignableScopes: pointer.To([]keyvault.RoleScope{
						scope,
					}),
				},
			}

			// TODO: @manicminer: when migrating to go-azure-sdk, the SDK should auto-retry on 409 responses and should consider manually polling afterwards

			if _, err = client.CreateOrUpdate(ctx, id.BaseURI(), id.Scope, id.RoleDefinitionName, payload); err != nil {
				return fmt.Errorf("creating %s: %v", id.ID(), err)
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context has no deadline")
			}
			stateConf := &pluginsdk.StateChangeConf{
				Pending: []string{"InProgress"},
				Target:  []string{"Found"},
				Refresh: func() (interface{}, string, error) {
					result, err := client.Get(ctx, id.BaseURI(), id.Scope, id.RoleDefinitionName)
					if err != nil {
						if response.WasNotFound(result.Response.Response) {
							return result, "InProgress", nil
						}

						return nil, "Error", err
					}

					return result, "Found", nil
				},
				ContinuousTargetOccurence: 5,
				PollInterval:              5 * time.Second,
				Timeout:                   time.Until(deadline),
			}
			if _, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for creation of %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r KeyVaultMHSMRoleDefinitionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedHSMs.DataPlaneRoleDefinitionsClient
			domainSuffix, ok := metadata.Client.Account.Environment.ManagedHSM.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Managed HSM domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			id, err := parse.ManagedHSMDataPlaneRoleDefinitionID(metadata.ResourceData.Id(), domainSuffix)
			if err != nil {
				return err
			}

			subscriptionId := commonids.NewSubscriptionID(metadata.Client.Account.SubscriptionId)
			managedHsmId, err := metadata.Client.ManagedHSMs.ManagedHSMIDFromBaseUrl(ctx, subscriptionId, id.BaseURI(), domainSuffix)
			if err != nil {
				return fmt.Errorf("determining the Managed HSM ID from the Base URI %q: %+v", id.BaseURI(), err)
			}
			if managedHsmId == nil {
				return fmt.Errorf("unable to determine the Managed HSM ID from the Base URI %q: %+v", id.BaseURI(), err)
			}

			locks.ByName(managedHsmId.ID(), "azurerm_key_vault_managed_hardware_security_module")
			defer locks.UnlockByName(managedHsmId.ID(), "azurerm_key_vault_managed_hardware_security_module")

			result, err := client.Get(ctx, id.BaseURI(), id.Scope, id.RoleDefinitionName)
			if err != nil {
				if response.WasNotFound(result.Response.Response) {
					return metadata.MarkAsGone(id)
				}
				return err
			}

			state := KeyVaultMHSMRoleDefinitionModel{
				Name:         pointer.From(result.Name),
				ManagedHSMID: managedHsmId.ID(),

				// TODO: remove in 4.0
				VaultBaseUrl: id.BaseURI(),
			}

			if v := pointer.From(result.ID); v != "" {
				roleID, err := roledefinitions.ParseScopedRoleDefinitionIDInsensitively(v)
				if err != nil {
					return fmt.Errorf("paring role definition id %q: %+v", v, err)
				}
				state.ResourceManagerId = roleID.ID()
			}

			if prop := result.RoleDefinitionProperties; prop != nil {
				state.Description = pointer.ToString(prop.Description)
				state.RoleType = string(prop.RoleType)
				state.RoleName = pointer.From(prop.RoleName)
				state.Permission = flattenKeyVaultMHSMRolePermission(prop.Permissions)
			}

			metadata.SetID(id)
			return metadata.Encode(&state)
		},
	}
}

func (r KeyVaultMHSMRoleDefinitionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: time.Minute * 10,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) (err error) {
			client := metadata.Client.ManagedHSMs.DataPlaneRoleDefinitionsClient
			domainSuffix, ok := metadata.Client.Account.Environment.ManagedHSM.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Managed HSM domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			id, err := parse.ManagedHSMDataPlaneRoleDefinitionID(metadata.ResourceData.Id(), domainSuffix)
			if err != nil {
				return err
			}

			subscriptionId := commonids.NewSubscriptionID(metadata.Client.Account.SubscriptionId)
			managedHsmId, err := metadata.Client.ManagedHSMs.ManagedHSMIDFromBaseUrl(ctx, subscriptionId, id.BaseURI(), domainSuffix)
			if err != nil {
				return fmt.Errorf("determining the Managed HSM ID from the Base URI %q: %+v", id.BaseURI(), err)
			}
			if managedHsmId == nil {
				return fmt.Errorf("unable to determine the Managed HSM ID from the Base URI %q: %+v", id.BaseURI(), err)
			}

			locks.ByName(managedHsmId.ID(), "azurerm_key_vault_managed_hardware_security_module")
			defer locks.UnlockByName(managedHsmId.ID(), "azurerm_key_vault_managed_hardware_security_module")

			var model KeyVaultMHSMRoleDefinitionModel
			if err = metadata.Decode(&model); err != nil {
				return err
			}

			existing, err := client.Get(ctx, id.BaseURI(), id.Scope, id.RoleDefinitionName)
			if err != nil {
				if response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("not found resource to update: %s", id)
				}
				return fmt.Errorf("retrieving role definition by name %s: %v", model.Name, err)
			}

			payload := keyvault.RoleDefinitionCreateParameters{
				Properties: &keyvault.RoleDefinitionProperties{
					RoleName:    pointer.To(model.RoleName),
					Description: pointer.To(model.Description),
					RoleType:    keyvault.RoleTypeCustomRole,
					Permissions: expandKeyVaultMHSMRolePermissions(model.Permission),
				},
			}

			_, err = client.CreateOrUpdate(ctx, id.BaseURI(), id.Scope, id.RoleDefinitionName, payload)
			if err != nil {
				return fmt.Errorf("updating %s: %v", id.ID(), err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r KeyVaultMHSMRoleDefinitionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedHSMs.DataPlaneRoleDefinitionsClient
			domainSuffix, ok := metadata.Client.Account.Environment.ManagedHSM.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Managed HSM domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			id, err := parse.ManagedHSMDataPlaneRoleDefinitionID(metadata.ResourceData.Id(), domainSuffix)
			if err != nil {
				return err
			}

			subscriptionId := commonids.NewSubscriptionID(metadata.Client.Account.SubscriptionId)
			managedHsmId, err := metadata.Client.ManagedHSMs.ManagedHSMIDFromBaseUrl(ctx, subscriptionId, id.BaseURI(), domainSuffix)
			if err != nil {
				return fmt.Errorf("determining the Managed HSM ID from the Base URI %q: %+v", id.BaseURI(), err)
			}
			if managedHsmId == nil {
				return fmt.Errorf("unable to determine the Managed HSM ID from the Base URI %q: %+v", id.BaseURI(), err)
			}

			locks.ByName(managedHsmId.ID(), "azurerm_key_vault_managed_hardware_security_module")
			defer locks.UnlockByName(managedHsmId.ID(), "azurerm_key_vault_managed_hardware_security_module")

			// TODO: @manicminer: when migrating to go-azure-sdk, the SDK should auto-retry on 409 responses
			// (these occur when a recently deleted assignment for the role has not yet fully replicated)

			if _, err = client.Delete(ctx, id.BaseURI(), id.Scope, id.RoleDefinitionName); err != nil {
				return fmt.Errorf("deleting %+v: %v", id, err)
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context has no deadline")
			}
			stateConf := &pluginsdk.StateChangeConf{
				Pending: []string{"InProgress"},
				Target:  []string{"NotFound"},
				Refresh: func() (interface{}, string, error) {
					result, err := client.Get(ctx, id.BaseURI(), id.Scope, id.RoleDefinitionName)
					if err != nil {
						if response.WasNotFound(result.Response.Response) {
							return result, "NotFound", nil
						}

						return nil, "Error", err
					}

					return result, "InProgress", nil
				},
				ContinuousTargetOccurence: 5,
				PollInterval:              5 * time.Second,
				Timeout:                   time.Until(deadline),
			}
			if _, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r KeyVaultMHSMRoleDefinitionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ManagedHSMDataPlaneRoleDefinitionID
}

func expandKeyVaultMHSMRolePermissions(perms []Permission) *[]keyvault.Permission {
	var res []keyvault.Permission
	for _, perm := range perms {
		var dataActions, notDataActions = make([]keyvault.DataAction, 0), make([]keyvault.DataAction, 0)
		for _, data := range perm.DataActions {
			dataActions = append(dataActions, keyvault.DataAction(data))
		}
		for _, notData := range perm.NotDataActions {
			notDataActions = append(notDataActions, keyvault.DataAction(notData))
		}

		res = append(res, keyvault.Permission{
			Actions:        pointer.To(perm.Actions),
			NotActions:     pointer.To(perm.NotActions),
			DataActions:    pointer.To(dataActions),
			NotDataActions: pointer.To(notDataActions),
		})
	}
	return &res
}

func flattenKeyVaultMHSMRolePermission(perms *[]keyvault.Permission) []Permission {
	if perms == nil {
		return make([]Permission, 0)
	}

	var res []Permission
	for _, perm := range *perms {
		var data, notData []string
		for _, item := range pointer.From(perm.DataActions) {
			data = append(data, string(item))
		}
		for _, item := range pointer.From(perm.NotDataActions) {
			notData = append(notData, string(item))
		}

		res = append(res, Permission{
			Actions:        pointer.From(perm.Actions),
			NotActions:     pointer.From(perm.NotActions),
			DataActions:    data,
			NotDataActions: notData,
		})
	}
	return res
}
