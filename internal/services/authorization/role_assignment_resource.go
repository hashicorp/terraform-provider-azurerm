// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package authorization

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2020-04-01-preview/authorization" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-05-01-preview/roledefinitions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-12-01/subscriptions"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	billingValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/billing/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// TODO: this wants splitting into virtual resources with Virtual IDs

func resourceArmRoleAssignment() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmRoleAssignmentCreate,
		Read:   resourceArmRoleAssignmentRead,
		Delete: resourceArmRoleAssignmentDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.RoleAssignmentID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"scope": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					// Elevated access (aka User Access Administrator role) is needed to assign roles in the following scopes:
					// https://docs.microsoft.com/en-us/azure/role-based-access-control/elevate-access-global-admin#azure-cli
					validation.StringMatch(regexp.MustCompile("/"), "Root scope (/) is invalid"),
					validation.StringMatch(regexp.MustCompile("/providers/Microsoft.Subscription.*"), "Subscription scope is invalid"),
					validation.StringMatch(regexp.MustCompile("/providers/Microsoft.Capacity"), "Capacity scope is invalid"),
					validation.StringMatch(regexp.MustCompile("/providers/Microsoft.BillingBenefits"), "BillingBenefits scope is invalid"),

					billingValidate.EnrollmentID,
					commonids.ValidateManagementGroupID,
					commonids.ValidateSubscriptionID,
					commonids.ValidateResourceGroupID,
					azure.ValidateResourceID,
				),
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

			"principal_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"principal_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"User",
					"Group",
					"ServicePrincipal",
				}, false),
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

			"condition": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"condition_version"},
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"condition_version": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"condition"},
				ValidateFunc: validation.StringInSlice([]string{
					"1.0",
					"2.0",
				}, false),
			},
		},
	}
}

func resourceArmRoleAssignmentCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	roleAssignmentsClient := meta.(*clients.Client).Authorization.RoleAssignmentsClient
	roleDefinitionsClient := meta.(*clients.Client).Authorization.ScopedRoleDefinitionsClient
	subscriptionClient := meta.(*clients.Client).Subscription.SubscriptionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	scope := d.Get("scope").(string)
	scopeId, err := commonids.ParseScopeID(scope)
	if err != nil {
		return fmt.Errorf("parsing %s: %+v", scopeId, err)
	}

	var roleDefinitionId string
	if v, ok := d.GetOk("role_definition_id"); ok {
		roleDefinitionId = v.(string)
	} else if v, ok := d.GetOk("role_definition_name"); ok {
		roleName := v.(string)
		roleDefinitions, err := roleDefinitionsClient.List(ctx, *scopeId, roledefinitions.ListOperationOptions{
			Filter: pointer.To(fmt.Sprintf("roleName eq '%s'", roleName)),
		})
		if err != nil {
			return fmt.Errorf("loading Role Definition List: %+v", err)
		}
		if roleDefinitions.Model == nil || len(*roleDefinitions.Model) != 1 {
			return fmt.Errorf("loading Role Definition List: could not find role '%s'", roleName)
		}
		roleDefinitionId = *(*roleDefinitions.Model)[0].Id
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

	existing, err := roleAssignmentsClient.Get(ctx, scope, name, tenantId)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing Role Assignment ID for %q (Scope %q): %+v", name, scope, err)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_role_assignment", *existing.ID)
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

	condition := d.Get("condition").(string)
	conditionVersion := d.Get("condition_version").(string)

	if condition != "" && conditionVersion != "" {
		properties.RoleAssignmentProperties.Condition = utils.String(condition)
		properties.RoleAssignmentProperties.ConditionVersion = utils.String(conditionVersion)
	} else if condition != "" || conditionVersion != "" {
		return fmt.Errorf("`condition` and `conditionVersion` should be both set or unset")
	}

	skipPrincipalCheck := d.Get("skip_service_principal_aad_check").(bool)
	if skipPrincipalCheck {
		properties.RoleAssignmentProperties.PrincipalType = authorization.ServicePrincipal
	}

	principalType := d.Get("principal_type").(string)
	if principalType != "" {
		properties.RoleAssignmentProperties.PrincipalType = authorization.PrincipalType(principalType)
	}

	// LinkedAuthorizationFailed may occur in cross tenant setup because of replication lag.
	// Let's retry this error for cross tenant setup and when we are skipping principal check.
	retryLinkedAuthorizationFailedError := len(delegatedManagedIdentityResourceID) > 0 && skipPrincipalCheck

	if err := pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), retryRoleAssignmentsClient(d, scope, name, properties, meta, tenantId, retryLinkedAuthorizationFailedError)); err != nil {
		return err
	}

	read, err := roleAssignmentsClient.Get(ctx, scope, name, tenantId)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Role Assignment ID for %q (Scope %q)", name, scope)
	}

	d.SetId(parse.ConstructRoleAssignmentId(*read.ID, tenantId))
	return resourceArmRoleAssignmentRead(d, meta)
}

func resourceArmRoleAssignmentRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Authorization.RoleAssignmentsClient
	roleDefinitionsClient := meta.(*clients.Client).Authorization.ScopedRoleDefinitionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.RoleAssignmentID(d.Id())
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
		d.Set("scope", normalizeScopeValue(pointer.From(props.Scope)))
		d.Set("role_definition_id", props.RoleDefinitionID)
		d.Set("principal_id", props.PrincipalID)
		d.Set("principal_type", props.PrincipalType)
		d.Set("delegated_managed_identity_resource_id", props.DelegatedManagedIdentityResourceID)
		d.Set("description", props.Description)
		d.Set("condition", props.Condition)
		d.Set("condition_version", props.ConditionVersion)

		// allows for import when role name is used (also if the role name changes a plan will show a diff)
		if roleDefResourceId := props.RoleDefinitionID; roleDefResourceId != nil {
			// Workaround for https://github.com/hashicorp/pandora/issues/3257
			// The role definition id returned does not contain scope when the role definition was on tenant level (management group or tenant).
			// And adding tenant id as scope will cause 404 response, so just adding a slash to parse that.
			if strings.HasPrefix(*roleDefResourceId, "/providers") {
				roleDefResourceId = pointer.To(fmt.Sprintf("/%s", *roleDefResourceId))
			}
			parsedRoleDefId, err := roledefinitions.ParseScopedRoleDefinitionID(*roleDefResourceId)
			if err != nil {
				return fmt.Errorf("parsing %q: %+v", *roleDefResourceId, err)
			}
			roleResp, err := roleDefinitionsClient.Get(ctx, *parsedRoleDefId)
			if err != nil {
				return fmt.Errorf("loading Role Definition %q: %+v", *roleDefResourceId, err)
			}
			if roleResp.Model != nil && roleResp.Model.Properties != nil {
				d.Set("role_definition_name", pointer.From(roleResp.Model.Properties.RoleName))
			}

		}
	}

	return nil
}

func resourceArmRoleAssignmentDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Authorization.RoleAssignmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parseRoleAssignmentId(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.scope, id.name, id.tenantId)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return err
		}
	}

	return nil
}

func retryRoleAssignmentsClient(d *pluginsdk.ResourceData, scope string, name string, properties authorization.RoleAssignmentCreateParameters, meta interface{}, tenantId string, retryLinkedAuthorizationFailedError bool) func() *pluginsdk.RetryError {
	return func() *pluginsdk.RetryError {
		roleAssignmentsClient := meta.(*clients.Client).Authorization.RoleAssignmentsClient
		ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
		defer cancel()

		resp, err := roleAssignmentsClient.Create(ctx, scope, name, properties)
		if err != nil {
			switch {
			case utils.ResponseErrorIsRetryable(err):
				return pluginsdk.RetryableError(err)
			case utils.ResponseWasStatusCode(resp.Response, 400) && strings.Contains(err.Error(), "PrincipalNotFound"):
				// When waiting for service principal to become available
				return pluginsdk.RetryableError(err)
			case retryLinkedAuthorizationFailedError && utils.ResponseWasForbidden(resp.Response) && strings.Contains(err.Error(), "LinkedAuthorizationFailed"):
				return pluginsdk.RetryableError(err)
			default:
				return pluginsdk.NonRetryableError(err)
			}
		}

		if resp.ID == nil {
			return pluginsdk.NonRetryableError(fmt.Errorf("creation of Role Assignment %q did not return an id value", name))
		}

		stateConf := &pluginsdk.StateChangeConf{
			Pending: []string{
				"pending",
			},
			Target: []string{
				"ready",
			},
			Refresh:                   roleAssignmentCreateStateRefreshFunc(ctx, roleAssignmentsClient, *resp.ID, tenantId),
			MinTimeout:                5 * time.Second,
			ContinuousTargetOccurence: 5,
			Timeout:                   d.Timeout(pluginsdk.TimeoutCreate),
		}

		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return pluginsdk.NonRetryableError(fmt.Errorf("failed waiting for Role Assignment %q to finish replicating: %+v", name, err))
		}

		return nil
	}
}

type roleAssignmentId struct {
	scope    string
	name     string
	tenantId string
}

func parseRoleAssignmentId(input string) (*roleAssignmentId, error) {
	tenantId := ""
	segments := strings.Split(input, "|")
	if len(segments) == 2 {
		tenantId = segments[1]
		input = segments[0]
	}

	segments = strings.Split(input, "/providers/Microsoft.Authorization/roleAssignments/")
	if len(segments) != 2 {
		return nil, fmt.Errorf("Expected Role Assignment ID to be in the format `{scope}/providers/Microsoft.Authorization/roleAssignments/{name}` but got %q", input)
	}

	// /{scope}/providers/Microsoft.Authorization/roleAssignments/{roleAssignmentName}
	id := roleAssignmentId{
		scope:    strings.TrimPrefix(segments[0], "/"),
		name:     segments[1],
		tenantId: tenantId,
	}
	return &id, nil
}

func roleAssignmentCreateStateRefreshFunc(ctx context.Context, client *authorization.RoleAssignmentsClient, roleID string, tenantId string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.GetByID(ctx, roleID, tenantId)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return resp, "pending", nil
			}
			return resp, "failed", err
		}
		return resp, "ready", nil
	}
}

func getTenantIdBySubscriptionId(ctx context.Context, client *subscriptions.SubscriptionsClient, subscriptionId string) (string, error) {
	id := commonids.NewSubscriptionID(subscriptionId)
	resp, err := client.Get(ctx, id)
	if err != nil {
		return "", fmt.Errorf("retrieving %s: %+v", id, err)
	}
	tenantId := ""
	if model := resp.Model; model != nil && model.TenantId != nil {
		tenantId = *model.TenantId
	}

	if tenantId == "" {
		return "", fmt.Errorf("retrieving %s: `tenantId` was nil", id)
	}
	return tenantId, nil
}

func normalizeScopeValue(scope string) (result string) {
	if rg, err := commonids.ParseResourceGroupIDInsensitively(scope); err == nil {
		return rg.ID()
	}
	// only check part of IDs, there are may be other specific resource types, like storage account id
	// we may need append these parse logics below when needed
	return scope
}
