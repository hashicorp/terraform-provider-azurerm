package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ResourceTagsId struct {
	SubscriptionId string
	ResourceGroup  string
	Provider       string
	Namespace      string
	ResourceName   string
}

func NewResourceTagsID(subscriptionId, resourceGroup, provider, namespace, resourceName string) ResourceTagsId {
	return ResourceTagsId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Provider:       provider,
		Namespace:      namespace,
		ResourceName:   resourceName,
	}
}

func (id ResourceTagsId) String() string {
	segments := []string{
		fmt.Sprintf("Subscription ID %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
		fmt.Sprintf("Provider %q", id.Provider),
		fmt.Sprintf("Resource Namespace %q", id.Namespace),
		fmt.Sprintf("Resource Name %q", id.ResourceName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Resource Tags", segmentsStr)
}

func (id ResourceTagsId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/%s/%s/%s/providers/Microsoft.Resources/tags/default"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Provider, id.Namespace, id.ResourceName)
}

func (id ResourceTagsId) ParentResourceID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/%s/%s/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Provider, id.Namespace, id.ResourceName)
}

// ResourceTagsID parses a ResourceTags ID into an ResourceTagsId struct
func ResourceTagsID(input string) (*ResourceTagsId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}
	splitId := strings.Split(input, "/")

	resourceId := ResourceTagsId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
		Provider:       id.Provider,
		Namespace:      splitId[7],
		ResourceName:   splitId[8],
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.Provider == "" {
		return nil, fmt.Errorf("ID was missing the 'provider' element")
	}

	if resourceId.Namespace == "" {
		return nil, fmt.Errorf("ID was missing the 'namespace' element")
	}

	if resourceId.ResourceName == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceName' element")
	}

	return &resourceId, nil
}
