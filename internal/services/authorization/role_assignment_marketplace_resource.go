package authorization

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2020-04-01-preview/authorization" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const (
	MarketplaceScope = "/providers/Microsoft.Marketplace"
)

type RoleAssignmentMarketplaceModel struct {
	PrincipalID                        string `tfschema:"principal_id"`
	Name                               string `tfschema:"name"`
	RoleDefinitionID                   string `tfschema:"role_definition_id"`
	RoleDefinitionName                 string `tfschema:"role_definition_name"`
	SkipServicePrincipalAadCheck       bool   `tfschema:"skip_service_principal_aad_check"`
	DelegatedManagedIdentityResourceID string `tfschema:"delegated_managed_identity_resource_id"`
	Description                        string `tfschema:"description"`
	PrincipalType                      string `tfschema:"principal_type"`
}

type RoleAssignmentMarketplaceResource struct{}

var _ sdk.Resource = RoleAssignmentMarketplaceResource{}

func (r RoleAssignmentMarketplaceResource) ResourceType() string {
	return "azurerm_role_assignment_marketplace"
}

func (r RoleAssignmentMarketplaceResource) ModelObject() interface{} {
	return &RoleAssignmentMarketplaceModel{}
}

func (r RoleAssignmentMarketplaceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.RoleAssignmentMarketplaceID
}

func (r RoleAssignmentMarketplaceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"principal_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"role_definition_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ExactlyOneOf: []string{"role_definition_name"},
			ValidateFunc: validate.RoleResourceID,
		},

		"role_definition_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ExactlyOneOf: []string{"role_definition_id"},
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"skip_service_principal_aad_check": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"delegated_managed_identity_resource_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r RoleAssignmentMarketplaceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"principal_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r RoleAssignmentMarketplaceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model RoleAssignmentMarketplaceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			roleAssignmentsClient := metadata.Client.Authorization.RoleAssignmentsClient
			roleDefinitionsClient := metadata.Client.Authorization.RoleDefinitionsClient
			subscriptionClient := metadata.Client.Subscription.Client
			subscriptionId := metadata.Client.Account.SubscriptionId
			roleDefinitionId := model.RoleDefinitionID

			if model.RoleDefinitionName != "" {
				roleDefinitions, err := roleDefinitionsClient.List(ctx, MarketplaceScope, fmt.Sprintf("roleName eq '%s'", model.RoleDefinitionName))
				if err != nil {
					return fmt.Errorf("loading Role Definition List: %+v", err)
				}

				if len(roleDefinitions.Values()) != 1 {
					return fmt.Errorf("loading Role Definition List: could not find role '%s'", model.RoleDefinitionName)
				}

				roleDefinitionId = *roleDefinitions.Values()[0].ID
			}

			name := model.Name
			if name == "" {
				uuid, err := uuid.GenerateUUID()
				if err != nil {
					return fmt.Errorf("generating UUID for Role Assignment: %+v", err)
				}

				name = uuid
			}

			tenantId := ""
			if model.DelegatedManagedIdentityResourceID != "" {
				var err error
				tenantId, err = getTenantIdBySubscriptionId(ctx, subscriptionClient, subscriptionId)
				if err != nil {
					return err
				}
			}

			id := parse.NewRoleAssignmentMarketplaceID(name, tenantId)
			existing, err := roleAssignmentsClient.Get(ctx, MarketplaceScope, id.Name, id.TenantId)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for presence of existing Role Assignment ID for %q: %+v", name, err)
				}
			}

			if existing.ID != nil && *existing.ID != "" {
				return tf.ImportAsExistsError("azurerm_role_assignment_marketplace", *existing.ID)
			}

			properties := authorization.RoleAssignmentCreateParameters{
				RoleAssignmentProperties: &authorization.RoleAssignmentProperties{
					RoleDefinitionID: &roleDefinitionId,
					PrincipalID:      &model.PrincipalID,
					Description:      &model.Description,
				},
			}

			if model.DelegatedManagedIdentityResourceID != "" {
				properties.RoleAssignmentProperties.DelegatedManagedIdentityResourceID = &model.DelegatedManagedIdentityResourceID
			}

			if model.SkipServicePrincipalAadCheck {
				properties.RoleAssignmentProperties.PrincipalType = authorization.ServicePrincipal
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("could not retrieve context deadline for %q", name)
			}

			if err = pluginsdk.Retry(time.Until(deadline), retryRoleAssignmentsClient(metadata.ResourceData, MarketplaceScope, name, properties, metadata.Client, tenantId)); err != nil {
				return err
			}

			read, err := roleAssignmentsClient.Get(ctx, MarketplaceScope, name, tenantId)
			if err != nil {
				return err
			}
			if read.ID == nil {
				return fmt.Errorf("cannot read Role Assignment ID for %q", name)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r RoleAssignmentMarketplaceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Authorization.RoleAssignmentsClient
			roleDefinitionsClient := metadata.Client.Authorization.RoleDefinitionsClient
			id, err := parse.RoleAssignmentMarketplaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.GetByID(ctx, id.AzureResourceID(), id.TenantId)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var model RoleAssignmentMarketplaceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			state := RoleAssignmentMarketplaceModel{
				Name:                         pointer.From(resp.Name),
				SkipServicePrincipalAadCheck: model.SkipServicePrincipalAadCheck,
			}

			if props := resp.RoleAssignmentPropertiesWithScope; props != nil {
				state.RoleDefinitionID = pointer.From(props.RoleDefinitionID)
				state.PrincipalID = pointer.From(props.PrincipalID)
				state.PrincipalType = string(props.PrincipalType)
				state.DelegatedManagedIdentityResourceID = pointer.From(props.DelegatedManagedIdentityResourceID)
				state.Description = pointer.From(props.Description)

				// allows for import when role name is used (also if the role name changes a plan will show a diff)
				if roleId := props.RoleDefinitionID; roleId != nil {
					roleResp, err := roleDefinitionsClient.GetByID(ctx, *roleId)
					if err != nil {
						return fmt.Errorf("loading Role Definition %q: %+v", *roleId, err)
					}

					if roleProps := roleResp.RoleDefinitionProperties; roleProps != nil {
						state.RoleDefinitionName = pointer.From(roleProps.RoleName)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r RoleAssignmentMarketplaceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Authorization.RoleAssignmentsClient
			id, err := parse.RoleAssignmentMarketplaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Delete(ctx, MarketplaceScope, id.Name, id.TenantId)
			if err != nil {
				if !utils.ResponseWasNotFound(resp.Response) {
					return err
				}
			}

			return nil
		},
	}
}
