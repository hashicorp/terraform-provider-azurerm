// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type RulesEngineId struct {
	SubscriptionId string
	ResourceGroup  string
	FrontdoorName  string
	Name           string
}

func NewRulesEngineID(subscriptionId, resourceGroup, frontdoorName, name string) RulesEngineId {
	return RulesEngineId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		FrontdoorName:  frontdoorName,
		Name:           name,
	}
}

func (id RulesEngineId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Frontdoor Name %q", id.FrontdoorName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Rules Engine", segmentsStr)
}

func (id RulesEngineId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/frontdoors/%s/rulesEngines/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.FrontdoorName, id.Name)
}

// RulesEngineID parses a RulesEngine ID into an RulesEngineId struct
func RulesEngineID(input string) (*RulesEngineId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an RulesEngine ID: %+v", input, err)
	}

	resourceId := RulesEngineId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.FrontdoorName, err = id.PopSegment("frontdoors"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("rulesEngines"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// RulesEngineIDInsensitively parses an RulesEngine ID into an RulesEngineId struct, insensitively
// This should only be used to parse an ID for rewriting, the RulesEngineID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func RulesEngineIDInsensitively(input string) (*RulesEngineId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := RulesEngineId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'frontdoors' segment
	frontdoorsKey := "frontdoors"
	for key := range id.Path {
		if strings.EqualFold(key, frontdoorsKey) {
			frontdoorsKey = key
			break
		}
	}
	if resourceId.FrontdoorName, err = id.PopSegment(frontdoorsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'rulesEngines' segment
	rulesEnginesKey := "rulesEngines"
	for key := range id.Path {
		if strings.EqualFold(key, rulesEnginesKey) {
			rulesEnginesKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(rulesEnginesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
