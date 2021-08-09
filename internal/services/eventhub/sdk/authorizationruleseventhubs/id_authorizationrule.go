package authorizationruleseventhubs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AuthorizationRuleId struct {
	SubscriptionId string
	ResourceGroup  string
	NamespaceName  string
	EventhubName   string
	Name           string
}

func NewAuthorizationRuleID(subscriptionId, resourceGroup, namespaceName, eventhubName, name string) AuthorizationRuleId {
	return AuthorizationRuleId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		NamespaceName:  namespaceName,
		EventhubName:   eventhubName,
		Name:           name,
	}
}

func (id AuthorizationRuleId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Eventhub Name %q", id.EventhubName),
		fmt.Sprintf("Namespace Name %q", id.NamespaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Authorization Rule", segmentsStr)
}

func (id AuthorizationRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventHub/namespaces/%s/eventhubs/%s/authorizationRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.NamespaceName, id.EventhubName, id.Name)
}

// AuthorizationRuleID parses a AuthorizationRule ID into an AuthorizationRuleId struct
func AuthorizationRuleID(input string) (*AuthorizationRuleId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AuthorizationRuleId{
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
	if resourceId.Name, err = id.PopSegment("authorizationRules"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// AuthorizationRuleIDInsensitively parses an AuthorizationRule ID into an AuthorizationRuleId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the AuthorizationRuleID method should be used instead for validation etc.
func AuthorizationRuleIDInsensitively(input string) (*AuthorizationRuleId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AuthorizationRuleId{
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

	// find the correct casing for the 'authorizationRules' segment
	authorizationRulesKey := "authorizationRules"
	for key := range id.Path {
		if strings.EqualFold(key, authorizationRulesKey) {
			authorizationRulesKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(authorizationRulesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
