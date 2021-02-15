package parse

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type EventSubscriptionId struct {
	Scope string
	Name  string
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
