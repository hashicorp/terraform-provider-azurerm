package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ApplicationId struct {
	SubscriptionId string
	ResourceGroup  string
	IoTAppName     string
}

func NewApplicationID(subscriptionId, resourceGroup, ioTAppName string) ApplicationId {
	return ApplicationId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		IoTAppName:     ioTAppName,
	}
}

func (id ApplicationId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.IoTCentral/IoTApps/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.IoTAppName)
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

	if resourceId.IoTAppName, err = id.PopSegment("IoTApps"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
