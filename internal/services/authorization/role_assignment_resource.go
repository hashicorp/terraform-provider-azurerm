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

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-04-01/roleassignments"
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
				ExactlyOneOf:     []string{"role_definition_id", "role_definition_name"},
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"role_definition_name": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ExactlyOneOf:     []string{"role_definition_name", "role_definition_id"},
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
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"condition_version": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"1.0",
					"2.0",
				}, false),
			},
		},
	}
}

func resourceArmRoleAssignmentCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	roleAssignmentsClient := meta.(*clients.Client).Authorization.ScopedRoleAssignmentsClient
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
	}

	if v, ok := d.GetOk("role_definition_name"); ok {
		roleName := v.(string)
		roleDefinitions, err := roleDefinitionsClient.List(ctx, commonids.NewScopeID(scope), roledefinitions.ListOperationOptions{
			Filter: pointer.To(fmt.Sprintf("roleName eq '%s'", roleName)),
		})
		if err != nil {
			return fmt.Errorf("listing role definitions: %+v", err)
		}
		if roleDefinitions.Model == nil || len(*roleDefinitions.Model) != 1 {
			return fmt.Errorf("listing role definitions: could not find role `%s`", roleName)
		}
		roleDefinitionId = *(*roleDefinitions.Model)[0].Id
	}
	d.Set("role_definition_id", roleDefinitionId)

	principalId := d.Get("principal_id").(string)

	if name == "" {
		generatedUUID, err := uuid.GenerateUUID()
		if err != nil {
			return fmt.Errorf("generating UUID for Role Assignment: %+v", err)
		}

		name = generatedUUID
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

	id := parse.NewScopedRoleAssignmentID(scope, name, tenantId)
	options := roleassignments.DefaultGetOperationOptions()
	if tenantId != "" {
		options.TenantId = pointer.To(tenantId)
	}

	existing, err := roleAssignmentsClient.Get(ctx, id.ScopedId, options)
	if err != nil && !response.WasNotFound(existing.HttpResponse) {
		return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_role_assignment", id.ID())
	}

	params := roleassignments.RoleAssignmentCreateParameters{
		Properties: roleassignments.RoleAssignmentProperties{
			RoleDefinitionId: roleDefinitionId,
			PrincipalId:      principalId,
			Description:      pointer.To(d.Get("description").(string)),
		},
	}
	props := &params.Properties

	if len(delegatedManagedIdentityResourceID) > 0 {
		props.DelegatedManagedIdentityResourceId = pointer.To(delegatedManagedIdentityResourceID)
	}

	condition := d.Get("condition").(string)
	conditionVersion := d.Get("condition_version").(string)

	switch {
	case condition != "" && conditionVersion != "":
		props.Condition = pointer.To(condition)
		props.ConditionVersion = pointer.To(conditionVersion)
	case condition != "":
		props.Condition = pointer.To(condition)
		props.ConditionVersion = pointer.To("2.0")
	case conditionVersion != "":
		return fmt.Errorf("`condition_version` should not be set without `condition`")
	}

	skipPrincipalCheck := d.Get("skip_service_principal_aad_check").(bool)
	if skipPrincipalCheck {
		props.PrincipalType = pointer.To(roleassignments.PrincipalTypeServicePrincipal)
	}

	if principalType := d.Get("principal_type").(string); principalType != "" {
		props.PrincipalType = pointer.To(roleassignments.PrincipalType(principalType))
	}

	// LinkedAuthorizationFailed may occur in cross tenant setup because of replication lag.
	// Let's retry this error for cross tenant setup and when we are skipping principal check.
	retryLinkedAuthorizationFailedError := len(delegatedManagedIdentityResourceID) > 0 && skipPrincipalCheck
	if err := pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), retryRoleAssignmentsClient(d, id, params, meta, retryLinkedAuthorizationFailedError)); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceArmRoleAssignmentRead(d, meta)
}

func resourceArmRoleAssignmentRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Authorization.ScopedRoleAssignmentsClient
	roleDefinitionsClient := meta.(*clients.Client).Authorization.ScopedRoleDefinitionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ScopedRoleAssignmentID(d.Id())
	if err != nil {
		return err
	}

	options := roleassignments.DefaultGetOperationOptions()
	if id.TenantId != "" {
		options.TenantId = pointer.To(id.TenantId)
	}

	resp, err := client.Get(ctx, id.ScopedId, options)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if model := resp.Model; model != nil {
		d.Set("name", model.Name)

		if props := model.Properties; props != nil {
			d.Set("scope", normalizeScopeValue(pointer.From(props.Scope)))
			d.Set("role_definition_id", props.RoleDefinitionId)
			d.Set("principal_id", props.PrincipalId)
			d.Set("principal_type", pointer.From(props.PrincipalType))
			d.Set("delegated_managed_identity_resource_id", props.DelegatedManagedIdentityResourceId)
			d.Set("description", props.Description)
			d.Set("condition", props.Condition)
			d.Set("condition_version", props.ConditionVersion)

			if roleDefResourceId := props.RoleDefinitionId; roleDefResourceId != "" {
				// Workaround for https://github.com/hashicorp/pandora/issues/3257
				// The role definition id returned does not contain scope when the role definition was on tenant level (management group or tenant).
				// And adding tenant id as scope will cause 404 response, so just adding a slash to parse that.
				if strings.HasPrefix(roleDefResourceId, "/providers") {
					roleDefResourceId = fmt.Sprintf("/%s", roleDefResourceId)
				}
				parsedRoleDefId, err := roledefinitions.ParseScopedRoleDefinitionID(roleDefResourceId)
				if err != nil {
					return fmt.Errorf("parsing %q: %+v", roleDefResourceId, err)
				}
				roleResp, err := roleDefinitionsClient.Get(ctx, *parsedRoleDefId)
				if err != nil {
					return fmt.Errorf("retrieving Role Definition %q: %+v", roleDefResourceId, err)
				}
				if roleResp.Model != nil && roleResp.Model.Properties != nil {
					d.Set("role_definition_name", pointer.From(roleResp.Model.Properties.RoleName))
				}
			}
		}
	}

	return nil
}

func resourceArmRoleAssignmentDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Authorization.ScopedRoleAssignmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ScopedRoleAssignmentID(d.Id())
	if err != nil {
		return err
	}

	options := roleassignments.DefaultDeleteOperationOptions()
	if id.TenantId != "" {
		options.TenantId = pointer.To(id.TenantId)
	}

	resp, err := client.Delete(ctx, id.ScopedId, options)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return err
		}
	}

	return nil
}

func retryRoleAssignmentsClient(d *pluginsdk.ResourceData, id parse.ScopedRoleAssignmentId, param roleassignments.RoleAssignmentCreateParameters, meta interface{}, retryLinkedAuthorizationFailedError bool) func() *pluginsdk.RetryError {
	return func() *pluginsdk.RetryError {
		roleAssignmentsClient := meta.(*clients.Client).Authorization.ScopedRoleAssignmentsClient
		ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
		defer cancel()

		resp, err := roleAssignmentsClient.Create(ctx, id.ScopedId, param)
		if err != nil {
			switch {
			case utils.ResponseErrorIsRetryable(err):
				return pluginsdk.RetryableError(err)
			case response.WasStatusCode(resp.HttpResponse, 400) && strings.Contains(err.Error(), "PrincipalNotFound"):
				// When waiting for service principal to become available
				return pluginsdk.RetryableError(err)
			case retryLinkedAuthorizationFailedError && response.WasForbidden(resp.HttpResponse) && strings.Contains(err.Error(), "LinkedAuthorizationFailed"):
				return pluginsdk.RetryableError(err)
			default:
				return pluginsdk.NonRetryableError(err)
			}
		}

		if resp.Model == nil || resp.Model.Id == nil {
			return pluginsdk.NonRetryableError(fmt.Errorf("creation of %s did not return an id value", id))
		}

		stateConf := &pluginsdk.StateChangeConf{
			Pending: []string{
				"pending",
			},
			Target: []string{
				"ready",
			},
			Refresh:                   roleAssignmentCreateStateRefreshFunc(ctx, roleAssignmentsClient, id),
			MinTimeout:                5 * time.Second,
			ContinuousTargetOccurence: 5,
			Timeout:                   d.Timeout(pluginsdk.TimeoutCreate),
		}

		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return pluginsdk.NonRetryableError(fmt.Errorf("failed waiting for Role Assignment %s to finish replicating: %+v", id, err))
		}

		return nil
	}
}

func roleAssignmentCreateStateRefreshFunc(ctx context.Context, client *roleassignments.RoleAssignmentsClient, id parse.ScopedRoleAssignmentId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		options := roleassignments.DefaultGetOperationOptions()
		if id.TenantId != "" {
			options.TenantId = pointer.To(id.TenantId)
		}

		resp, err := client.Get(ctx, id.ScopedId, options)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
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
