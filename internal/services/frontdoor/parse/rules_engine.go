package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type RulesEngineId struct {
	SubscriptionId string
	ResourceGroup  string
	FrontdoorName  string
	Name           string
}

func NewRulesEngineID(subscriptionId, resourceGroup, frontdoorName, name string) RulesEngineId {
	return RulesEngineId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		FrontdoorName:  frontdoorName,
		Name:           name,
	}
}

func (id RulesEngineId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Frontdoor Name %q", id.FrontdoorName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Rules Engine", segmentsStr)
}

func (id RulesEngineId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/frontdoors/%s/rulesengines/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.FrontdoorName, id.Name)
}

// RulesEngineID parses a RulesEngine ID into an RulesEngineId struct
func RulesEngineID(input string) (*RulesEngineId, error) {
	id, err := azure.ParseAzureResourceID(input)
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

	if resourceId.FrontdoorName, err = id.PopSegment("frontdoors"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("rulesengines"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
