package parse

import (
	"fmt"
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
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
		managementGroupId, _ := ManagementGroupID(input) // if this is a management group ID, there should not be any error.
		scopeId.ManagementGroupId = managementGroupId.GroupId
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
	_, err := ManagementGroupID(input)
	return err == nil
}

// TODO -- move this to management group RP directory
type ManagementGroupId struct {
	GroupId string
}

func ManagementGroupID(input string) (*ManagementGroupId, error) {
	regex := regexp.MustCompile(`^/providers/[Mm]icrosoft\.[Mm]anagement/management[Gg]roups/`)
	if !regex.MatchString(input) {
		return nil, fmt.Errorf("unable to parse Management Group ID %q", input)
	}

	// Split the input ID by the regex
	segments := regex.Split(input, -1)
	if len(segments) != 2 {
		return nil, fmt.Errorf("unable to parse Management Group ID %q: expected id to have two segments after splitting", input)
	}
	groupID := segments[1]
	// portal says: The name can only be an ASCII letter, digit, -, _, (, ), . and have a maximum length constraint of 90
	if matched := regexp.MustCompile(`^[a-zA-Z0-9_().-]{1,90}$`).Match([]byte(groupID)); !matched {
		return nil, fmt.Errorf("unable to parse Management Group ID %q: group id can only consist of ASCII letters, digits, -, _, (, ), . , and cannot exceed the maximum length of 90", input)
	}

	id := ManagementGroupId{
		GroupId: groupID,
	}

	return &id, nil
}
