// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SpringCloudAppAssociationId struct {
	SubscriptionId string
	ResourceGroup  string
	SpringName     string
	AppName        string
	BindingName    string
}

func NewSpringCloudAppAssociationID(subscriptionId, resourceGroup, springName, appName, bindingName string) SpringCloudAppAssociationId {
	return SpringCloudAppAssociationId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		SpringName:     springName,
		AppName:        appName,
		BindingName:    bindingName,
	}
}

func (id SpringCloudAppAssociationId) String() string {
	segments := []string{
		fmt.Sprintf("Binding Name %q", id.BindingName),
		fmt.Sprintf("App Name %q", id.AppName),
		fmt.Sprintf("Spring Name %q", id.SpringName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Spring Cloud App Association", segmentsStr)
}

func (id SpringCloudAppAssociationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/apps/%s/bindings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SpringName, id.AppName, id.BindingName)
}

// SpringCloudAppAssociationID parses a SpringCloudAppAssociation ID into an SpringCloudAppAssociationId struct
func SpringCloudAppAssociationID(input string) (*SpringCloudAppAssociationId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SpringCloudAppAssociation ID: %+v", input, err)
	}

	resourceId := SpringCloudAppAssociationId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.SpringName, err = id.PopSegment("spring"); err != nil {
		return nil, err
	}
	if resourceId.AppName, err = id.PopSegment("apps"); err != nil {
		return nil, err
	}
	if resourceId.BindingName, err = id.PopSegment("bindings"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// SpringCloudAppAssociationIDInsensitively parses an SpringCloudAppAssociation ID into an SpringCloudAppAssociationId struct, insensitively
// This should only be used to parse an ID for rewriting, the SpringCloudAppAssociationID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func SpringCloudAppAssociationIDInsensitively(input string) (*SpringCloudAppAssociationId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SpringCloudAppAssociationId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'spring' segment
	springKey := "spring"
	for key := range id.Path {
		if strings.EqualFold(key, springKey) {
			springKey = key
			break
		}
	}
	if resourceId.SpringName, err = id.PopSegment(springKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'apps' segment
	appsKey := "apps"
	for key := range id.Path {
		if strings.EqualFold(key, appsKey) {
			appsKey = key
			break
		}
	}
	if resourceId.AppName, err = id.PopSegment(appsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'bindings' segment
	bindingsKey := "bindings"
	for key := range id.Path {
		if strings.EqualFold(key, bindingsKey) {
			bindingsKey = key
			break
		}
	}
	if resourceId.BindingName, err = id.PopSegment(bindingsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
