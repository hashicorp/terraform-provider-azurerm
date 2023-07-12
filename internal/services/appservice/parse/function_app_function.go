// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type FunctionAppFunctionId struct {
	SubscriptionId string
	ResourceGroup  string
	SiteName       string
	FunctionName   string
}

func NewFunctionAppFunctionID(subscriptionId, resourceGroup, siteName, functionName string) FunctionAppFunctionId {
	return FunctionAppFunctionId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		SiteName:       siteName,
		FunctionName:   functionName,
	}
}

func (id FunctionAppFunctionId) String() string {
	segments := []string{
		fmt.Sprintf("Function Name %q", id.FunctionName),
		fmt.Sprintf("Site Name %q", id.SiteName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Function App Function", segmentsStr)
}

func (id FunctionAppFunctionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/functions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SiteName, id.FunctionName)
}

// FunctionAppFunctionID parses a FunctionAppFunction ID into an FunctionAppFunctionId struct
func FunctionAppFunctionID(input string) (*FunctionAppFunctionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an FunctionAppFunction ID: %+v", input, err)
	}

	resourceId := FunctionAppFunctionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.SiteName, err = id.PopSegment("sites"); err != nil {
		return nil, err
	}
	if resourceId.FunctionName, err = id.PopSegment("functions"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
