package parse

import (
	"fmt"
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	parseMgmtGroup "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managementgroup/parse"
)

type RemediationId struct {
	Name string
	RemediationScopeId
}

type RemediationScopeId interface {
	ScopeId() string
}

type RemediationScopeAtSubscription struct {
	SubscriptionId string
	scopeId        string
}

func (id RemediationScopeAtSubscription) ScopeId() string {
	return id.scopeId
}

type RemediationScopeAtResourceGroup struct {
	scopeId        string
	SubscriptionId string
	ResourceGroup  string
}

func (id RemediationScopeAtResourceGroup) ScopeId() string {
	return id.scopeId
}

type RemediationScopeAtResource struct {
	scopeId string
}

func (id RemediationScopeAtResource) ScopeId() string {
	return id.scopeId
}

type RemediationScopeAtManagementGroup struct {
	scopeId           string
	ManagementGroupId string
}

func (id RemediationScopeAtManagementGroup) ScopeId() string {
	return id.scopeId
}

// TODO: This paring function is currently suppressing every case difference due to github issue: https://github.com/Azure/azure-rest-api-specs/issues/8353
// Currently the returned Remediation response from the service will have all the IDs converted into lower cases
func RemediationID(input string) (*RemediationId, error) {
	// in general, the id of a remediation should be:
	// {scope}/providers/Microsoft.PolicyInsights/remediations/{name}
	regex := regexp.MustCompile(`/providers/[Mm]icrosoft\.[Pp]olicy[Ii]nsights/remediations/`)
	if !regex.MatchString(input) {
		return nil, fmt.Errorf("unable to parse Policy Remediation ID %q", input)
	}

	segments := regex.Split(input, -1)

	if len(segments) != 2 {
		return nil, fmt.Errorf("unable to parse Policy Remediation ID %q: Expected 2 segments after split", input)
	}

	scope := segments[0]
	name := segments[1]
	if name == "" {
		return nil, fmt.Errorf("unable to parse Policy Remediation ID %q: remediation name is empty", input)
	}

	scopeId, err := RemediationScopeID(scope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Policy Remediation ID %q: %+v", input, err)
	}

	return &RemediationId{
		Name:               name,
		RemediationScopeId: scopeId,
	}, nil
}

func RemediationScopeID(input string) (RemediationScopeId, error) {
	if input == "" {
		return nil, fmt.Errorf("unable to parse Remediation Scope ID: ID is empty")
	}

	if isManagementGroupId(input) {
		id, _ := parseMgmtGroup.ManagementGroupID(input)
		return RemediationScopeAtManagementGroup{
			scopeId:           input,
			ManagementGroupId: id.GroupID,
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
		return RemediationScopeAtSubscription{
			scopeId:        input,
			SubscriptionId: id.SubscriptionID,
		}, nil
	}
	if id.ResourceGroup != "" && noExtraSegments {
		// should be a resourceGroup id
		return RemediationScopeAtResourceGroup{
			scopeId:        input,
			SubscriptionId: id.SubscriptionID,
			ResourceGroup:  id.ResourceGroup,
		}, nil
	}
	// should be a resource ID
	return RemediationScopeAtResource{
		scopeId: input,
	}, nil
}

func isManagementGroupId(input string) bool {
	_, err := parseMgmtGroup.ManagementGroupID(input)
	return err == nil
}
