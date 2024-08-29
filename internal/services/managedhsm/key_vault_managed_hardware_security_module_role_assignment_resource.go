// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedhsm

import (
	"context"
	"fmt"
	"regexp"
	"strings"
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

type KeyVaultManagedHSMRoleAssignmentModel struct {
	ManagedHSMID     string `tfschema:"managed_hsm_id"`
	Name             string `tfschema:"name"`
	Scope            string `tfschema:"scope"`
	RoleDefinitionId string `tfschema:"role_definition_id"`
	PrincipalId      string `tfschema:"principal_id"`
	ResourceId       string `tfschema:"resource_id"`

	// TODO: remove in v4.0
	VaultBaseUrl string `tfschema:"vault_base_url,removedInNextMajorVersion"`
}

var _ sdk.ResourceWithStateMigration = KeyVaultManagedHSMRoleAssignmentResource{}

type KeyVaultManagedHSMRoleAssignmentResource struct{}

func (r KeyVaultManagedHSMRoleAssignmentResource) Arguments() map[string]*pluginsdk.Schema {
	s := map[string]*pluginsdk.Schema{
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

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"scope": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile(`^(/|/keys|/keys/.+)$`), "scope should be one of `/`, `/keys', `/keys/<key_name>`"),
		},

		"role_definition_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: roledefinitions.ValidateScopedRoleDefinitionID,
		},

		"principal_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
	if !features.FourPointOhBeta() {
		s["vault_base_url"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			Deprecated:   "The field `vault_base_url` has been deprecated in favour of `managed_hsm_id` and will be removed in 4.0 of the Azure Provider",
			ValidateFunc: validation.StringIsNotEmpty,
		}
	}
	return s
}

func (r KeyVaultManagedHSMRoleAssignmentResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r KeyVaultManagedHSMRoleAssignmentResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.ManagedHSMRoleAssignmentV0ToV1{},
		},
	}
}

func (r KeyVaultManagedHSMRoleAssignmentResource) ModelObject() interface{} {
	return &KeyVaultManagedHSMRoleAssignmentModel{}
}

func (r KeyVaultManagedHSMRoleAssignmentResource) ResourceType() string {
	return "azurerm_key_vault_managed_hardware_security_module_role_assignment"
}

func (r KeyVaultManagedHSMRoleAssignmentResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedHSMs.DataPlaneRoleAssignmentsClient
			domainSuffix, ok := metadata.Client.Account.Environment.ManagedHSM.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Managed HSM domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			var config KeyVaultManagedHSMRoleAssignmentModel
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
					return fmt.Errorf("parsing the Data Plane Endpoint %q: %+v", pointer.From(endpoint), err)
				}
			}

			if managedHsmId == nil && !features.FourPointOhBeta() {
				endpoint, err = parse.ManagedHSMEndpoint(config.VaultBaseUrl, domainSuffix)
				if err != nil {
					return fmt.Errorf("parsing the Data Plane Endpoint %q: %+v", pointer.From(endpoint), err)
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

			locks.ByName(managedHsmId.ID(), "azurerm_key_vault_managed_hardware_security_module")
			defer locks.UnlockByName(managedHsmId.ID(), "azurerm_key_vault_managed_hardware_security_module")

			id := parse.NewManagedHSMDataPlaneRoleAssignmentID(endpoint.ManagedHSMName, endpoint.DomainSuffix, config.Scope, config.Name)
			existing, err := client.Get(ctx, endpoint.BaseURI(), config.Scope, config.Name)
			if !utils.ResponseWasNotFound(existing.Response) {
				if err != nil {
					return fmt.Errorf("retrieving %s: %v", id.ID(), err)
				}
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			var param keyvault.RoleAssignmentCreateParameters
			param.Properties = &keyvault.RoleAssignmentProperties{
				PrincipalID: pointer.FromString(config.PrincipalId),
				// the role definition id may have '/' prefix, but the api doesn't accept it
				RoleDefinitionID: pointer.FromString(strings.TrimPrefix(config.RoleDefinitionId, "/")),
			}

			//nolint:misspell
			// TODO: @manicminer: when migrating to go-azure-sdk, the SDK should auto-retry on 400 responses with code "BadParameter" and message "Unkown role definition" (note the misspelling)

			if _, err = client.Create(ctx, endpoint.BaseURI(), config.Scope, config.Name, param); err != nil {
				return fmt.Errorf("creating %s: %v", id.ID(), err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r KeyVaultManagedHSMRoleAssignmentResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedHSMs.DataPlaneRoleAssignmentsClient
			domainSuffix, ok := metadata.Client.Account.Environment.ManagedHSM.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Managed HSM domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			id, err := parse.ManagedHSMDataPlaneRoleAssignmentID(metadata.ResourceData.Id(), domainSuffix)
			if err != nil {
				return err
			}

			subscriptionId := commonids.NewSubscriptionID(metadata.Client.Account.SubscriptionId)
			resourceManagerId, err := metadata.Client.ManagedHSMs.ManagedHSMIDFromBaseUrl(ctx, subscriptionId, id.BaseURI(), domainSuffix)
			if err != nil {
				return fmt.Errorf("determining Resource Manager ID for %q: %+v", id, err)
			}
			if resourceManagerId == nil {
				return fmt.Errorf("unable to determine the Resource Manager ID for %s", id)
			}

			resp, err := client.Get(ctx, id.BaseURI(), id.Scope, id.RoleAssignmentName)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			model := KeyVaultManagedHSMRoleAssignmentModel{
				ManagedHSMID: resourceManagerId.ID(),
				Name:         id.RoleAssignmentName,
				Scope:        id.Scope,
			}

			if !features.FourPointOhBeta() {
				model.VaultBaseUrl = id.BaseURI()
			}

			if props := resp.Properties; props != nil {
				model.PrincipalId = pointer.From(props.PrincipalID)
				model.ResourceId = pointer.From(resp.ID) // TODO: verify if we should deprecate this

				if roleDefinitionId := pointer.From(props.RoleDefinitionID); roleDefinitionId != "" {
					parsed, err := roledefinitions.ParseScopedRoleDefinitionIDInsensitively(roleDefinitionId)
					if err != nil {
						return fmt.Errorf("parsing role definition id %q: %v", roleDefinitionId, err)
					}

					model.RoleDefinitionId = parsed.ID()
				}
			}

			return metadata.Encode(&model)
		},
	}
}

func (r KeyVaultManagedHSMRoleAssignmentResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedHSMs.DataPlaneRoleAssignmentsClient

			domainSuffix, ok := metadata.Client.Account.Environment.ManagedHSM.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Managed HSM domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			id, err := parse.ManagedHSMDataPlaneRoleAssignmentID(metadata.ResourceData.Id(), domainSuffix)
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

			if _, err := client.Delete(ctx, id.BaseURI(), id.Scope, id.RoleAssignmentName); err != nil {
				return fmt.Errorf("deleting %s: %v", id.ID(), err)
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context has no deadline")
			}
			stateConf := &pluginsdk.StateChangeConf{
				Pending: []string{"InProgress"},
				Target:  []string{"NotFound"},
				Refresh: func() (interface{}, string, error) {
					result, err := client.Get(ctx, id.BaseURI(), id.Scope, id.RoleAssignmentName)
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

func (r KeyVaultManagedHSMRoleAssignmentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ManagedHSMDataPlaneRoleAssignmentID
}
