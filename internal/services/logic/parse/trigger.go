package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type TriggerId struct {
	SubscriptionId string
	ResourceGroup  string
	WorkflowName   string
	Name           string
}

func NewTriggerID(subscriptionId, resourceGroup, workflowName, name string) TriggerId {
	return TriggerId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		WorkflowName:   workflowName,
		Name:           name,
	}
}

func (id TriggerId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Workflow Name %q", id.WorkflowName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Trigger", segmentsStr)
}

func (id TriggerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logic/workflows/%s/triggers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WorkflowName, id.Name)
}

// TriggerID parses a Trigger ID into an TriggerId struct
func TriggerID(input string) (*TriggerId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := TriggerId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.WorkflowName, err = id.PopSegment("workflows"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("triggers"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
