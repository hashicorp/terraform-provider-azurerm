package policy

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/policyinsights/mgmt/2019-10-01/policyinsights"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy/parse"
)

// RemediationCreateUpdateAtScope is a wrapper of the 4 CreateOrUpdate functions on RemediationsClient, combining them into one to simplify code.
func RemediationCreateUpdateAtScope(ctx context.Context, client *policyinsights.RemediationsClient, name string, scope parse.RemediationScopeId, parameters policyinsights.Remediation) (policyinsights.Remediation, error) {
	switch scope.Type {
	case parse.AtSubscription:
		return client.CreateOrUpdateAtSubscription(ctx, scope.SubscriptionId, name, parameters)
	case parse.AtResourceGroup:
		return client.CreateOrUpdateAtResourceGroup(ctx, scope.SubscriptionId, scope.ResourceGroup, name, parameters)
	case parse.AtResource:
		return client.CreateOrUpdateAtResource(ctx, scope.ScopeId, name, parameters)
	case parse.AtManagementGroup:
		return client.CreateOrUpdateAtManagementGroup(ctx, scope.ManagementGroupId, name, parameters)
	default:
		return policyinsights.Remediation{}, fmt.Errorf("Error creating Policy Remediation: Invalid scope type %q", scope.Type)
	}
}

// RemediationGetAtScope is a wrapper of the 4 Get functions on RemediationsClient, combining them into one to simplify code.
func RemediationGetAtScope(ctx context.Context, client *policyinsights.RemediationsClient, name string, scope parse.RemediationScopeId) (policyinsights.Remediation, error) {
	switch scope.Type {
	case parse.AtSubscription:
		return client.GetAtSubscription(ctx, scope.SubscriptionId, name)
	case parse.AtResourceGroup:
		return client.GetAtResourceGroup(ctx, scope.SubscriptionId, scope.ResourceGroup, name)
	case parse.AtResource:
		return client.GetAtResource(ctx, scope.ScopeId, name)
	case parse.AtManagementGroup:
		return client.GetAtManagementGroup(ctx, scope.ManagementGroupId, name)
	default:
		return policyinsights.Remediation{}, fmt.Errorf("Error reading Policy Remediation: Invalid scope type %q", scope.Type)
	}
}

// RemediationDeleteAtScope is a wrapper of the 4 Delete functions on RemediationsClient, combining them into one to simplify code.
func RemediationDeleteAtScope(ctx context.Context, client *policyinsights.RemediationsClient, name string, scope parse.RemediationScopeId) (policyinsights.Remediation, error) {
	switch scope.Type {
	case parse.AtSubscription:
		return client.DeleteAtSubscription(ctx, scope.SubscriptionId, name)
	case parse.AtResourceGroup:
		return client.DeleteAtResourceGroup(ctx, scope.SubscriptionId, scope.ResourceGroup, name)
	case parse.AtResource:
		return client.DeleteAtResource(ctx, scope.ScopeId, name)
	case parse.AtManagementGroup:
		return client.DeleteAtManagementGroup(ctx, scope.ManagementGroupId, name)
	default:
		return policyinsights.Remediation{}, fmt.Errorf("Error deleting Policy Remediation: Invalid scope type %q", scope.Type)
	}
}
