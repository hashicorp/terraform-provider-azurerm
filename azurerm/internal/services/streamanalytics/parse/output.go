package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type OutputId struct {
	SubscriptionId   string
	ResourceGroup    string
	StreamingjobName string
	Name             string
}

func NewOutputID(subscriptionId, resourceGroup, streamingjobName, name string) OutputId {
	return OutputId{
		SubscriptionId:   subscriptionId,
		ResourceGroup:    resourceGroup,
		StreamingjobName: streamingjobName,
		Name:             name,
	}
}

func (id OutputId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Streamingjob Name %q", id.StreamingjobName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Output", segmentsStr)
}

func (id OutputId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StreamAnalytics/streamingjobs/%s/outputs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.StreamingjobName, id.Name)
}

// OutputID parses a Output ID into an OutputId struct
func OutputID(input string) (*OutputId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := OutputId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.StreamingjobName, err = id.PopSegment("streamingjobs"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("outputs"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
