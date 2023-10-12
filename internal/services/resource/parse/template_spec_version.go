// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type TemplateSpecVersionId struct {
	SubscriptionId   string
	ResourceGroup    string
	TemplateSpecName string
	VersionName      string
}

func NewTemplateSpecVersionID(subscriptionId, resourceGroup, templateSpecName, versionName string) TemplateSpecVersionId {
	return TemplateSpecVersionId{
		SubscriptionId:   subscriptionId,
		ResourceGroup:    resourceGroup,
		TemplateSpecName: templateSpecName,
		VersionName:      versionName,
	}
}

func (id TemplateSpecVersionId) String() string {
	segments := []string{
		fmt.Sprintf("Version Name %q", id.VersionName),
		fmt.Sprintf("Template Spec Name %q", id.TemplateSpecName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Template Spec Version", segmentsStr)
}

func (id TemplateSpecVersionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Resources/templateSpecs/%s/versions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.TemplateSpecName, id.VersionName)
}

// TemplateSpecVersionID parses a TemplateSpecVersion ID into an TemplateSpecVersionId struct
func TemplateSpecVersionID(input string) (*TemplateSpecVersionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an TemplateSpecVersion ID: %+v", input, err)
	}

	resourceId := TemplateSpecVersionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.TemplateSpecName, err = id.PopSegment("templateSpecs"); err != nil {
		return nil, err
	}
	if resourceId.VersionName, err = id.PopSegment("versions"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
