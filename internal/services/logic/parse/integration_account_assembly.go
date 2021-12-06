package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type IntegrationAccountAssemblyId struct {
	SubscriptionId         string
	ResourceGroup          string
	IntegrationAccountName string
	AssemblyName           string
}

func NewIntegrationAccountAssemblyID(subscriptionId, resourceGroup, integrationAccountName, assemblyName string) IntegrationAccountAssemblyId {
	return IntegrationAccountAssemblyId{
		SubscriptionId:         subscriptionId,
		ResourceGroup:          resourceGroup,
		IntegrationAccountName: integrationAccountName,
		AssemblyName:           assemblyName,
	}
}

func (id IntegrationAccountAssemblyId) String() string {
	segments := []string{
		fmt.Sprintf("Assembly Name %q", id.AssemblyName),
		fmt.Sprintf("Integration Account Name %q", id.IntegrationAccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Integration Account Assembly", segmentsStr)
}

func (id IntegrationAccountAssemblyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logic/integrationAccounts/%s/assemblies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.IntegrationAccountName, id.AssemblyName)
}

// IntegrationAccountAssemblyID parses a IntegrationAccountAssembly ID into an IntegrationAccountAssemblyId struct
func IntegrationAccountAssemblyID(input string) (*IntegrationAccountAssemblyId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := IntegrationAccountAssemblyId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.IntegrationAccountName, err = id.PopSegment("integrationAccounts"); err != nil {
		return nil, err
	}
	if resourceId.AssemblyName, err = id.PopSegment("assemblies"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
