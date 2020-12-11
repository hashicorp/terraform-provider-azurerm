package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ApplicationId struct {
	SubscriptionId   string
	ResourceGroup    string
	BatchAccountName string
	Name             string
}

func NewApplicationID(subscriptionId, resourceGroup, batchAccountName, name string) ApplicationId {
	return ApplicationId{
		SubscriptionId:   subscriptionId,
		ResourceGroup:    resourceGroup,
		BatchAccountName: batchAccountName,
		Name:             name,
	}
}

func (id ApplicationId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Batch Account Name %q", id.BatchAccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Application", segmentsStr)
}

func (id ApplicationId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Batch/batchAccounts/%s/applications/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.BatchAccountName, id.Name)
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

	if resourceId.BatchAccountName, err = id.PopSegment("batchAccounts"); err != nil {
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
