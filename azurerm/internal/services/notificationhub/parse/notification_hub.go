package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type NotificationHubId struct {
	ResourceGroup string
	NamespaceName string
	Name          string
}

func NotificationHubID(input string) (*NotificationHubId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Notification Hub ID %q: %+v", input, err)
	}

	app := NotificationHubId{
		ResourceGroup: id.ResourceGroup,
	}

	if app.NamespaceName, err = id.PopSegment("namespaces"); err != nil {
		return nil, err
	}

	if app.Name, err = id.PopSegment("notificationHubs"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &app, nil
}
