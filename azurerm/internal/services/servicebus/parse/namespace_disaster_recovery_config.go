package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type NamespaceDisasterRecoveryConfigId struct {
	SubscriptionId             string
	ResourceGroup              string
	NamespaceName              string
	DisasterRecoveryConfigName string
}

func NewNamespaceDisasterRecoveryConfigID(subscriptionId, resourceGroup, namespaceName, disasterRecoveryConfigName string) NamespaceDisasterRecoveryConfigId {
	return NamespaceDisasterRecoveryConfigId{
		SubscriptionId:             subscriptionId,
		ResourceGroup:              resourceGroup,
		NamespaceName:              namespaceName,
		DisasterRecoveryConfigName: disasterRecoveryConfigName,
	}
}

func (id NamespaceDisasterRecoveryConfigId) String() string {
	segments := []string{
		fmt.Sprintf("Disaster Recovery Config Name %q", id.DisasterRecoveryConfigName),
		fmt.Sprintf("Namespace Name %q", id.NamespaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Namespace Disaster Recovery Config", segmentsStr)
}

func (id NamespaceDisasterRecoveryConfigId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ServiceBus/namespaces/%s/disasterRecoveryConfigs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.NamespaceName, id.DisasterRecoveryConfigName)
}

// NamespaceDisasterRecoveryConfigID parses a NamespaceDisasterRecoveryConfig ID into an NamespaceDisasterRecoveryConfigId struct
func NamespaceDisasterRecoveryConfigID(input string) (*NamespaceDisasterRecoveryConfigId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := NamespaceDisasterRecoveryConfigId{
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
	if resourceId.DisasterRecoveryConfigName, err = id.PopSegment("disasterRecoveryConfigs"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
