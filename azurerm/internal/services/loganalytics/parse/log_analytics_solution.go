package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LogAnalyticsSolutionId struct {
	SubscriptionId string
	ResourceGroup  string
	SolutionName   string
}

func NewLogAnalyticsSolutionID(subscriptionId, resourceGroup, solutionName string) LogAnalyticsSolutionId {
	return LogAnalyticsSolutionId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		SolutionName:   solutionName,
	}
}

func (id LogAnalyticsSolutionId) String() string {
	segments := []string{
		fmt.Sprintf("Solution Name %q", id.SolutionName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Log Analytics Solution", segmentsStr)
}

func (id LogAnalyticsSolutionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationsManagement/solutions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SolutionName)
}

// LogAnalyticsSolutionID parses a LogAnalyticsSolution ID into an LogAnalyticsSolutionId struct
func LogAnalyticsSolutionID(input string) (*LogAnalyticsSolutionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := LogAnalyticsSolutionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.SolutionName, err = id.PopSegment("solutions"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
