package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type PublicMaintenanceConfigurationId struct {
	SubscriptionId string
	Name           string
}

func NewPublicMaintenanceConfigurationID(subscriptionId, name string) PublicMaintenanceConfigurationId {
	return PublicMaintenanceConfigurationId{
		SubscriptionId: subscriptionId,
		Name:           name,
	}
}

func (id PublicMaintenanceConfigurationId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Public Maintenance Configuration", segmentsStr)
}

func (id PublicMaintenanceConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Maintenance/publicMaintenanceConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.Name)
}

// PublicMaintenanceConfigurationID parses a PublicMaintenanceConfiguration ID into an PublicMaintenanceConfigurationId struct
func PublicMaintenanceConfigurationID(input string) (*PublicMaintenanceConfigurationId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := PublicMaintenanceConfigurationId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.Name, err = id.PopSegment("publicMaintenanceConfigurations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
