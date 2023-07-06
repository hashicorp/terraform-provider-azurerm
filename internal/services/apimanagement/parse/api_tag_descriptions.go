// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ApiTagDescriptionsId struct {
	SubscriptionId     string
	ResourceGroup      string
	ServiceName        string
	ApiName            string
	TagDescriptionName string
}

func NewApiTagDescriptionsID(subscriptionId, resourceGroup, serviceName, apiName, tagDescriptionName string) ApiTagDescriptionsId {
	return ApiTagDescriptionsId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		ServiceName:        serviceName,
		ApiName:            apiName,
		TagDescriptionName: tagDescriptionName,
	}
}

func (id ApiTagDescriptionsId) String() string {
	segments := []string{
		fmt.Sprintf("Tag Description Name %q", id.TagDescriptionName),
		fmt.Sprintf("Api Name %q", id.ApiName),
		fmt.Sprintf("Service Name %q", id.ServiceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Api Tag Descriptions", segmentsStr)
}

func (id ApiTagDescriptionsId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/apis/%s/tagDescriptions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServiceName, id.ApiName, id.TagDescriptionName)
}

// ApiTagDescriptionsID parses a ApiTagDescriptions ID into an ApiTagDescriptionsId struct
func ApiTagDescriptionsID(input string) (*ApiTagDescriptionsId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ApiTagDescriptions ID: %+v", input, err)
	}

	resourceId := ApiTagDescriptionsId{
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
	if resourceId.ApiName, err = id.PopSegment("apis"); err != nil {
		return nil, err
	}
	if resourceId.TagDescriptionName, err = id.PopSegment("tagDescriptions"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
