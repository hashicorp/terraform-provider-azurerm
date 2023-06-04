package authorization

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2020-04-01-preview/authorization" // nolint: staticcheck
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const (
	MarketplaceScope = "/providers/Microsoft.Marketplace"
)

func resourceArmRoleAssignmentMarketplace() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmRoleAssignmentMarketplaceCreate,
		Read:   resourceArmRoleAssignmentMarketplaceRead,
		Delete: resourceArmRoleAssignmentMarketplaceDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.RoleAssignmentMarketplaceID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"principal_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"role_definition_id": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ConflictsWith:    []string{"role_definition_name"},
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"role_definition_name": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ConflictsWith:    []string{"role_definition_id"},
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validation.StringIsNotEmpty,
			},

			"principal_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"skip_service_principal_aad_check": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Computed: true,
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
		},
	}
}

func resourceArmRoleAssignmentMarketplaceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	roleAssignmentsClient := meta.(*clients.Client).Authorization.RoleAssignmentsClient
	roleDefinitionsClient := meta.(*clients.Client).Authorization.RoleDefinitionsClient
	subscriptionClient := meta.(*clients.Client).Subscription.Client
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)

	var roleDefinitionId string
	if v, ok := d.GetOk("role_definition_id"); ok {
		roleDefinitionId = v.(string)
	} else if v, ok := d.GetOk("role_definition_name"); ok {
		roleName := v.(string)
		roleDefinitions, err := roleDefinitionsClient.List(ctx, MarketplaceScope, fmt.Sprintf("roleName eq '%s'", roleName))
		if err != nil {
			return fmt.Errorf("loading Role Definition List: %+v", err)
		}
		if len(roleDefinitions.Values()) != 1 {
			return fmt.Errorf("loading Role Definition List: could not find role '%s'", roleName)
		}
		roleDefinitionId = *roleDefinitions.Values()[0].ID
	} else {
		return fmt.Errorf("Error: either role_definition_id or role_definition_name needs to be set")
	}
	d.Set("role_definition_id", roleDefinitionId)

	principalId := d.Get("principal_id").(string)

	if name == "" {
		uuid, err := uuid.GenerateUUID()
		if err != nil {
			return fmt.Errorf("generating UUID for Role Assignment: %+v", err)
		}

		name = uuid
	}

	tenantId := ""
	delegatedManagedIdentityResourceID := d.Get("delegated_managed_identity_resource_id").(string)
	if len(delegatedManagedIdentityResourceID) > 0 {
		var err error
		tenantId, err = getTenantIdBySubscriptionId(ctx, subscriptionClient, subscriptionId)
		if err != nil {
			return err
		}
	}

	existing, err := roleAssignmentsClient.Get(ctx, MarketplaceScope, name, tenantId)
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
			RoleDefinitionID: utils.String(roleDefinitionId),
			PrincipalID:      utils.String(principalId),
			Description:      utils.String(d.Get("description").(string)),
		},
	}

	if len(delegatedManagedIdentityResourceID) > 0 {
		properties.RoleAssignmentProperties.DelegatedManagedIdentityResourceID = utils.String(delegatedManagedIdentityResourceID)
	}

	skipPrincipalCheck := d.Get("skip_service_principal_aad_check").(bool)
	if skipPrincipalCheck {
		properties.RoleAssignmentProperties.PrincipalType = authorization.ServicePrincipal
	}

	if err := pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), retryRoleAssignmentsClient(d, MarketplaceScope, name, properties, meta, tenantId)); err != nil {
		return err
	}

	read, err := roleAssignmentsClient.Get(ctx, MarketplaceScope, name, tenantId)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("cannot read Role Assignment ID for %q", name)
	}

	d.SetId(parse.ConstructRoleAssignmentId(*read.ID, tenantId))
	return resourceArmRoleAssignmentMarketplaceRead(d, meta)
}

func resourceArmRoleAssignmentMarketplaceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Authorization.RoleAssignmentsClient
	roleDefinitionsClient := meta.(*clients.Client).Authorization.RoleDefinitionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.RoleAssignmentMarketplaceID(d.Id())
	if err != nil {
		return err
	}
	resp, err := client.GetByID(ctx, id.AzureResourceID(), id.TenantId)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Role Assignment ID %q was not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("loading Role Assignment %q: %+v", d.Id(), err)
	}

	d.Set("name", resp.Name)

	if props := resp.RoleAssignmentPropertiesWithScope; props != nil {
		d.Set("role_definition_id", props.RoleDefinitionID)
		d.Set("principal_id", props.PrincipalID)
		d.Set("principal_type", props.PrincipalType)
		d.Set("delegated_managed_identity_resource_id", props.DelegatedManagedIdentityResourceID)
		d.Set("description", props.Description)

		// allows for import when role name is used (also if the role name changes a plan will show a diff)
		if roleId := props.RoleDefinitionID; roleId != nil {
			roleResp, err := roleDefinitionsClient.GetByID(ctx, *roleId)
			if err != nil {
				return fmt.Errorf("loading Role Definition %q: %+v", *roleId, err)
			}

			if roleProps := roleResp.RoleDefinitionProperties; roleProps != nil {
				d.Set("role_definition_name", roleProps.RoleName)
			}
		}
	}

	return nil
}

func resourceArmRoleAssignmentMarketplaceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Authorization.RoleAssignmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.RoleAssignmentMarketplaceID(d.Id())
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
}
