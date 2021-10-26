package frontdoors

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type RulesEngineId struct {
	SubscriptionId string
	ResourceGroup  string
	FrontDoorName  string
	Name           string
}

func NewRulesEngineID(subscriptionId, resourceGroup, frontDoorName, name string) RulesEngineId {
	return RulesEngineId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		FrontDoorName:  frontDoorName,
		Name:           name,
	}
}

func (id RulesEngineId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Front Door Name %q", id.FrontDoorName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Rules Engine", segmentsStr)
}

func (id RulesEngineId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/frontDoors/%s/rulesEngines/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.FrontDoorName, id.Name)
}

// ParseRulesEngineID parses a RulesEngine ID into an RulesEngineId struct
func ParseRulesEngineID(input string) (*RulesEngineId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := RulesEngineId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.FrontDoorName, err = id.PopSegment("frontDoors"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("rulesEngines"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParseRulesEngineIDInsensitively parses an RulesEngine ID into an RulesEngineId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseRulesEngineID method should be used instead for validation etc.
func ParseRulesEngineIDInsensitively(input string) (*RulesEngineId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := RulesEngineId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'frontDoors' segment
	frontDoorsKey := "frontDoors"
	for key := range id.Path {
		if strings.EqualFold(key, frontDoorsKey) {
			frontDoorsKey = key
			break
		}
	}
	if resourceId.FrontDoorName, err = id.PopSegment(frontDoorsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'rulesEngines' segment
	rulesEnginesKey := "rulesEngines"
	for key := range id.Path {
		if strings.EqualFold(key, rulesEnginesKey) {
			rulesEnginesKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(rulesEnginesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
