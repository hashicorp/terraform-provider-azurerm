package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ApplicationId struct {
	SubscriptionId       string
	ResourceGroup        string
	ApplicationGroupName string
	Name                 string
}

func NewApplicationID(subscriptionId, resourceGroup, applicationGroupName, name string) ApplicationId {
	return ApplicationId{
		SubscriptionId:       subscriptionId,
		ResourceGroup:        resourceGroup,
		ApplicationGroupName: applicationGroupName,
		Name:                 name,
	}
}

func (id ApplicationId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Application Group Name %q", id.ApplicationGroupName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Application", segmentsStr)
}

func (id ApplicationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DesktopVirtualization/applicationGroups/%s/applications/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ApplicationGroupName, id.Name)
}

// ApplicationID parses a Application ID into an ApplicationId struct
func ApplicationID(input string) (*ApplicationId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ApplicationId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ApplicationGroupName, err = id.PopSegment("applicationGroups"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("applications"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ApplicationIDInsensitively parses an Application ID into an ApplicationId struct, insensitively
// This should only be used to parse an ID for rewriting, the ApplicationID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func ApplicationIDInsensitively(input string) (*ApplicationId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ApplicationId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'applicationGroups' segment
	applicationGroupsKey := "applicationGroups"
	for key := range id.Path {
		if strings.EqualFold(key, applicationGroupsKey) {
			applicationGroupsKey = key
			break
		}
	}
	if resourceId.ApplicationGroupName, err = id.PopSegment(applicationGroupsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'applications' segment
	applicationsKey := "applications"
	for key := range id.Path {
		if strings.EqualFold(key, applicationsKey) {
			applicationsKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(applicationsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
