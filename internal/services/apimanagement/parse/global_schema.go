// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type GlobalSchemaId struct {
	SubscriptionId string
	ResourceGroup  string
	ServiceName    string
	SchemaName     string
}

func NewGlobalSchemaID(subscriptionId, resourceGroup, serviceName, schemaName string) GlobalSchemaId {
	return GlobalSchemaId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ServiceName:    serviceName,
		SchemaName:     schemaName,
	}
}

func (id GlobalSchemaId) String() string {
	segments := []string{
		fmt.Sprintf("Schema Name %q", id.SchemaName),
		fmt.Sprintf("Service Name %q", id.ServiceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Global Schema", segmentsStr)
}

func (id GlobalSchemaId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/schemas/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServiceName, id.SchemaName)
}

// GlobalSchemaID parses a GlobalSchema ID into an GlobalSchemaId struct
func GlobalSchemaID(input string) (*GlobalSchemaId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an GlobalSchema ID: %+v", input, err)
	}

	resourceId := GlobalSchemaId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ServiceName, err = id.PopSegment("service"); err != nil {
		return nil, err
	}
	if resourceId.SchemaName, err = id.PopSegment("schemas"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
