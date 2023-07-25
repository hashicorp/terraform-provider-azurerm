// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

var _ resourceids.Id = ResourceProviderId{}

type ResourceProviderId struct {
	SubscriptionId   string
	ResourceProvider string
}

func NewResourceProviderID(subscriptionId, resourceProvider string) ResourceProviderId {
	return ResourceProviderId{
		SubscriptionId:   subscriptionId,
		ResourceProvider: resourceProvider,
	}
}

func (id ResourceProviderId) ID() string {
	fmtString := "/subscriptions/%s/providers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceProvider)
}

func (id ResourceProviderId) String() string {
	return fmt.Sprintf("Resource Provider %q", id.ResourceProvider)
}

// ResourceProviderID parses a ResourceProvider ID into an ResourceProviderId struct
func ResourceProviderID(input string) (*ResourceProviderId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ResourceProviderId{
		SubscriptionId:   id.SubscriptionID,
		ResourceProvider: id.Provider,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceProvider == "" {
		return nil, fmt.Errorf("ID was missing the 'providers' element")
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
