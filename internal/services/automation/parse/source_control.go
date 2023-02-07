package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SourceControlId struct {
	SubscriptionId        string
	ResourceGroup         string
	AutomationAccountName string
	Name                  string
}

func NewSourceControlID(subscriptionId, resourceGroup, automationAccountName, name string) SourceControlId {
	return SourceControlId{
		SubscriptionId:        subscriptionId,
		ResourceGroup:         resourceGroup,
		AutomationAccountName: automationAccountName,
		Name:                  name,
	}
}

func (id SourceControlId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Automation Account Name %q", id.AutomationAccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Source Control", segmentsStr)
}

func (id SourceControlId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s/sourceControls/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.AutomationAccountName, id.Name)
}

// SourceControlID parses a SourceControl ID into an SourceControlId struct
func SourceControlID(input string) (*SourceControlId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SourceControlId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.AutomationAccountName, err = id.PopSegment("automationAccounts"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("sourceControls"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// SourceControlIDInsensitively parses an SourceControl ID into an SourceControlId struct, insensitively
// This should only be used to parse an ID for rewriting, the SourceControlID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func SourceControlIDInsensitively(input string) (*SourceControlId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SourceControlId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'automationAccounts' segment
	automationAccountsKey := "automationAccounts"
	for key := range id.Path {
		if strings.EqualFold(key, automationAccountsKey) {
			automationAccountsKey = key
			break
		}
	}
	if resourceId.AutomationAccountName, err = id.PopSegment(automationAccountsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'sourceControls' segment
	sourceControlsKey := "sourceControls"
	for key := range id.Path {
		if strings.EqualFold(key, sourceControlsKey) {
			sourceControlsKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(sourceControlsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
