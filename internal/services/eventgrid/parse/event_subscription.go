package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

var _ resourceids.Id = EventSubscriptionId{}

type EventSubscriptionId struct {
	Scope string
	Name  string
}

func NewEventSubscriptionID(scope string, name string) EventSubscriptionId {
	return EventSubscriptionId{
		Scope: scope,
		Name:  name,
	}
}

func (id EventSubscriptionId) ID() string {
	fmtStr := "%s/providers/Microsoft.EventGrid/eventSubscriptions/%s"
	return fmt.Sprintf(fmtStr, id.Scope, id.Name)
}

func (id EventSubscriptionId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Scope %q", id.Scope),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Event Subscription", segmentsStr)
}

func EventSubscriptionID(input string) (*EventSubscriptionId, error) {
	_, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse EventGrid Event Subscription ID %q: %+v", input, err)
	}

	segments := strings.Split(input, "/providers/Microsoft.EventGrid/eventSubscriptions/")
	if len(segments) != 2 {
		return nil, fmt.Errorf("Expected ID to be in the format `{scope}/providers/Microsoft.EventGrid/eventSubscriptions/{name} - got %d segments", len(segments))
	}

	eventSubscription := EventSubscriptionId{
		Scope: segments[0],
		Name:  segments[1],
	}

	return &eventSubscription, nil
}
