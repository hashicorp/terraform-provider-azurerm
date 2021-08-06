package consumergroups

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ConsumergroupId struct {
	SubscriptionId string
	ResourceGroup  string
	NamespaceName  string
	EventhubName   string
	Name           string
}

func NewConsumergroupID(subscriptionId, resourceGroup, namespaceName, eventhubName, name string) ConsumergroupId {
	return ConsumergroupId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		NamespaceName:  namespaceName,
		EventhubName:   eventhubName,
		Name:           name,
	}
}

func (id ConsumergroupId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Eventhub Name %q", id.EventhubName),
		fmt.Sprintf("Namespace Name %q", id.NamespaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Consumergroup", segmentsStr)
}

func (id ConsumergroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventHub/namespaces/%s/eventhubs/%s/consumergroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.NamespaceName, id.EventhubName, id.Name)
}

// ConsumergroupID parses a Consumergroup ID into an ConsumergroupId struct
func ConsumergroupID(input string) (*ConsumergroupId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ConsumergroupId{
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
	if resourceId.EventhubName, err = id.PopSegment("eventhubs"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("consumergroups"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ConsumergroupIDInsensitively parses an Consumergroup ID into an ConsumergroupId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ConsumergroupID method should be used instead for validation etc.
func ConsumergroupIDInsensitively(input string) (*ConsumergroupId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ConsumergroupId{
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

	// find the correct casing for the 'eventhubs' segment
	eventhubsKey := "eventhubs"
	for key := range id.Path {
		if strings.EqualFold(key, eventhubsKey) {
			eventhubsKey = key
			break
		}
	}
	if resourceId.EventhubName, err = id.PopSegment(eventhubsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'consumergroups' segment
	consumergroupsKey := "consumergroups"
	for key := range id.Path {
		if strings.EqualFold(key, consumergroupsKey) {
			consumergroupsKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(consumergroupsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
