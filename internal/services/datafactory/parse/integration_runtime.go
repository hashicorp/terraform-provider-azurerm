// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type IntegrationRuntimeId struct {
	SubscriptionId string
	ResourceGroup  string
	FactoryName    string
	Name           string
}

func NewIntegrationRuntimeID(subscriptionId, resourceGroup, factoryName, name string) IntegrationRuntimeId {
	return IntegrationRuntimeId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		FactoryName:    factoryName,
		Name:           name,
	}
}

func (id IntegrationRuntimeId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Factory Name %q", id.FactoryName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Integration Runtime", segmentsStr)
}

func (id IntegrationRuntimeId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataFactory/factories/%s/integrationruntimes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.FactoryName, id.Name)
}

// IntegrationRuntimeID parses a IntegrationRuntime ID into an IntegrationRuntimeId struct
func IntegrationRuntimeID(input string) (*IntegrationRuntimeId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an IntegrationRuntime ID: %+v", input, err)
	}

	resourceId := IntegrationRuntimeId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, errors.New("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, errors.New("ID was missing the 'resourceGroups' element")
	}

	if resourceId.FactoryName, err = id.PopSegment("factories"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("integrationruntimes"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
