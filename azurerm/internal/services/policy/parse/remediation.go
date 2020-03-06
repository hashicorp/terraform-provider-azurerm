package parse

import (
	"fmt"
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	parseMgmtGroup "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managementgroup/parse"
)

type ProvisioningType int

const (
	AtSubscription ProvisioningType = iota
	AtManagementGroup
	AtResourceGroup
	AtResource
)

type RemediationId struct {
	Name string
	RemediationScopeId
}

// TODO: This paring function is currently suppressing every case difference due to github issue: https://github.com/Azure/azure-rest-api-specs/issues/8353
// Currently the returned Remediation response from the service will have all the IDs converted into lower cases.
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
	scopeId, err := RemediationScopeID(scope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Policy Remediation ID %q: %+v", input, err)
	}

	name := segments[1]
	if name == "" {
		return nil, fmt.Errorf("unable to parse Policy Remediation ID %q: remediation name is empty", input)
	}

	id := RemediationId{
		Name:               name,
		RemediationScopeId: *scopeId,
	}

	return &id, nil
}

type RemediationScopeId struct {
	ScopeId           string
	Type              ProvisioningType
	SubscriptionId    string
	ResourceGroup     string
	ManagementGroupId string
}

func RemediationScopeID(input string) (*RemediationScopeId, error) {
	scopeId := RemediationScopeId{
		ScopeId: input,
	}

	if isManagementGroupId(input) {
		managementGroupId, _ := parseMgmtGroup.ManagementGroupID(input) // if this is a management group ID, there should not be any error.
		scopeId.ManagementGroupId = managementGroupId.GroupID
		scopeId.Type = AtManagementGroup
	} else {
		id, err := azure.ParseAzureResourceID(input)
		if err != nil {
			return nil, err
		}
		scopeId.SubscriptionId = id.SubscriptionID
		scopeId.ResourceGroup = id.ResourceGroup
		if id.ResourceGroup == "" {
			// it is a subscription id
			scopeId.Type = AtSubscription
		} else if err := id.ValidateNoEmptySegments(input); err == nil {
			// it is a resource group id
			scopeId.Type = AtResourceGroup
		} else {
			// it is a resource id
			scopeId.Type = AtResource
		}
	}

	return &scopeId, nil
}

func isManagementGroupId(input string) bool {
	_, err := parseMgmtGroup.ManagementGroupID(input)
	return err == nil
}
