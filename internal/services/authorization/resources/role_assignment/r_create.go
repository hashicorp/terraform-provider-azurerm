// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package role_assignment

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-04-01/roleassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-05-01-preview/roledefinitions"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmRoleAssignmentCreate(d *pluginsdk.ResourceData, meta any) error {
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
		tenantId, err = helpers.GetTenantIdBySubscriptionId(ctx, subscriptionClient, subscriptionId)
		if err != nil {
			return err
		}
	}

	id := parse.NewScopedRoleAssignmentID(scope, name, tenantId)
	options := roleassignments.DefaultGetOperationOptions()
	if tenantId != "" {
		options.TenantId = pointer.To(tenantId)
	}

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := roleAssignmentsClient.Get(ctx, id.ScopedId, options)
		if err != nil && !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_role_assignment", id.ID())
		}
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
		props.PrincipalType = pointer.ToEnum[roleassignments.PrincipalType](principalType)
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
