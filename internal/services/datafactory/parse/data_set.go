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

type DataSetId struct {
	SubscriptionId string
	ResourceGroup  string
	FactoryName    string
	Name           string
}

func NewDataSetID(subscriptionId, resourceGroup, factoryName, name string) DataSetId {
	return DataSetId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		FactoryName:    factoryName,
		Name:           name,
	}
}

func (id DataSetId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Factory Name %q", id.FactoryName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Data Set", segmentsStr)
}

func (id DataSetId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataFactory/factories/%s/datasets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.FactoryName, id.Name)
}

// DataSetID parses a DataSet ID into an DataSetId struct
func DataSetID(input string) (*DataSetId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an DataSet ID: %+v", input, err)
	}

	resourceId := DataSetId{
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
	if resourceId.Name, err = id.PopSegment("datasets"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
