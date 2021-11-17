package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type FeatureId struct {
	SubscriptionId    string
	ProviderNamespace string
	Name              string
}

func NewFeatureID(subscriptionId, resourceProvider, name string) FeatureId {
	return FeatureId{
		SubscriptionId:    subscriptionId,
		ProviderNamespace: resourceProvider,
		Name:              name,
	}
}

func (id FeatureId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("Feature: (%s)", segmentsStr)
}

func (id FeatureId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Features/providers/%s/features/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ProviderNamespace, id.Name)
}

// FeatureID parses a Feature ID into an FeatureId struct
func FeatureID(input string) (*FeatureId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FeatureId{
		SubscriptionId:    id.SubscriptionID,
		ProviderNamespace: id.SecondaryProvider,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ProviderNamespace == "" {
		return nil, fmt.Errorf("ID was missing the 'providers' element")
	}

	if resourceId.Name, err = id.PopSegment("features"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
