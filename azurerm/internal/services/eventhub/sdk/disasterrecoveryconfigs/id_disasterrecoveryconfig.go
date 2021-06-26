package disasterrecoveryconfigs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type DisasterRecoveryConfigId struct {
	SubscriptionId string
	ResourceGroup  string
	NamespaceName  string
	Name           string
}

func NewDisasterRecoveryConfigID(subscriptionId, resourceGroup, namespaceName, name string) DisasterRecoveryConfigId {
	return DisasterRecoveryConfigId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		NamespaceName:  namespaceName,
		Name:           name,
	}
}

func (id DisasterRecoveryConfigId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Namespace Name %q", id.NamespaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Disaster Recovery Config", segmentsStr)
}

func (id DisasterRecoveryConfigId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventHub/namespaces/%s/disasterRecoveryConfigs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.NamespaceName, id.Name)
}

// DisasterRecoveryConfigID parses a DisasterRecoveryConfig ID into an DisasterRecoveryConfigId struct
func DisasterRecoveryConfigID(input string) (*DisasterRecoveryConfigId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DisasterRecoveryConfigId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.NamespaceName, err = id.PopSegment("namespaces"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("disasterRecoveryConfigs"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// DisasterRecoveryConfigIDInsensitively parses an DisasterRecoveryConfig ID into an DisasterRecoveryConfigId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the DisasterRecoveryConfigID method should be used instead for validation etc.
func DisasterRecoveryConfigIDInsensitively(input string) (*DisasterRecoveryConfigId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DisasterRecoveryConfigId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'namespaces' segment
	namespacesKey := "namespaces"
	for key := range id.Path {
		if strings.EqualFold(key, namespacesKey) {
			namespacesKey = key
			break
		}
	}
	if resourceId.NamespaceName, err = id.PopSegment(namespacesKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'disasterRecoveryConfigs' segment
	disasterRecoveryConfigsKey := "disasterRecoveryConfigs"
	for key := range id.Path {
		if strings.EqualFold(key, disasterRecoveryConfigsKey) {
			disasterRecoveryConfigsKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(disasterRecoveryConfigsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
