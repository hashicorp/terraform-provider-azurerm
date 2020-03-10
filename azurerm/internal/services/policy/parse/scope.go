package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	parseMgmtGroup "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managementgroup/parse"
)

type PolicyScopeId interface {
	ScopeId() string
}

type ScopeAtSubscription struct {
	SubscriptionId string
	scopeId        string
}

func (id ScopeAtSubscription) ScopeId() string {
	return id.scopeId
}

type ScopeAtResourceGroup struct {
	scopeId        string
	SubscriptionId string
	ResourceGroup  string
}

func (id ScopeAtResourceGroup) ScopeId() string {
	return id.scopeId
}

type ScopeAtResource struct {
	scopeId string
}

func (id ScopeAtResource) ScopeId() string {
	return id.scopeId
}

type ScopeAtManagementGroup struct {
	scopeId           string
	ManagementGroupId string
}

func (id ScopeAtManagementGroup) ScopeId() string {
	return id.scopeId
}
func PolicyScopeID(input string) (PolicyScopeId, error) {
	if input == "" {
		return nil, fmt.Errorf("unable to parse Remediation Scope ID: ID is empty")
	}

	if isManagementGroupId(input) {
		id, _ := parseMgmtGroup.ManagementGroupID(input)
		return ScopeAtManagementGroup{
			scopeId:           input,
			ManagementGroupId: id.GroupId,
		}, nil
	}
	// scope is not a management group ID, should be subscription ID, resource group ID or a resource ID
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Policy Remediation ID %q: %+v", input, err)
	}
	noExtraSegments := id.ValidateNoEmptySegments(input) == nil
	if id.ResourceGroup == "" && noExtraSegments {
		// should be a subscription id
		return ScopeAtSubscription{
			scopeId:        input,
			SubscriptionId: id.SubscriptionID,
		}, nil
	}
	if id.ResourceGroup != "" && noExtraSegments {
		// should be a resourceGroup id
		return ScopeAtResourceGroup{
			scopeId:        input,
			SubscriptionId: id.SubscriptionID,
			ResourceGroup:  id.ResourceGroup,
		}, nil
	}
	// should be a resource ID
	return ScopeAtResource{
		scopeId: input,
	}, nil
}

func isManagementGroupId(input string) bool {
	_, err := parseMgmtGroup.ManagementGroupID(input)
	return err == nil
}
