package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SecretValueId struct {
	SubscriptionId string
	ResourceGroup  string
	SecretName     string
	ValueName      string
}

func NewSecretValueID(subscriptionId, resourceGroup, secretName, valueName string) SecretValueId {
	return SecretValueId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		SecretName:     secretName,
		ValueName:      valueName,
	}
}

func (id SecretValueId) String() string {
	segments := []string{
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
		fmt.Sprintf("Secret Name %q", id.SecretName),
		fmt.Sprintf("Value Name %q", id.ValueName),
	}
	return strings.Join(segments, " / ")
}

func (id SecretValueId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ServiceFabricMesh/secrets/%s/values/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SecretName, id.ValueName)
}

// SecretValueID parses a SecretValue ID into an SecretValueId struct
func SecretValueID(input string) (*SecretValueId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SecretValueId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.SecretName, err = id.PopSegment("secrets"); err != nil {
		return nil, err
	}
	if resourceId.ValueName, err = id.PopSegment("values"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
