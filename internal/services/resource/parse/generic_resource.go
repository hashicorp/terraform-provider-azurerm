package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// it's only to parse the existing resource.

type GenericResourceId struct {
	SubscriptionId   string
	ResourceProvider string
	ResourceType     string
}

// GenericResourceID parses a Generic Resource ID into an GenericResourceId struct
func GenericResourceID(input string) (*GenericResourceId, error) {
	id, err := resourceids.ParseAzureResourceID(input) //azure.ParseAzureResourceID(input)

	if err != nil {
		return nil, err
	}

	resourceTypeBuilder := strings.Builder{}
	addSlash := false
	for k, _ := range id.Path {
		if addSlash {
			resourceTypeBuilder.WriteString("/")
		}
		resourceTypeBuilder.WriteString(k)
		addSlash = true
	}

	resourceId := GenericResourceId{
		SubscriptionId:   id.SubscriptionID,
		ResourceProvider: id.Provider,
		ResourceType:     resourceTypeBuilder.String(),
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceProvider == "" {
		return nil, fmt.Errorf("ID was missing the 'providers' element")
	}

	if resourceId.ResourceType == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceType' element")
	}

	return &resourceId, nil
}
