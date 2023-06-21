// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type HostPoolRegistrationInfoId struct {
	SubscriptionId       string
	ResourceGroup        string
	HostPoolName         string
	RegistrationInfoName string
}

func NewHostPoolRegistrationInfoID(subscriptionId, resourceGroup, hostPoolName, registrationInfoName string) HostPoolRegistrationInfoId {
	return HostPoolRegistrationInfoId{
		SubscriptionId:       subscriptionId,
		ResourceGroup:        resourceGroup,
		HostPoolName:         hostPoolName,
		RegistrationInfoName: registrationInfoName,
	}
}

func (id HostPoolRegistrationInfoId) String() string {
	segments := []string{
		fmt.Sprintf("Registration Info Name %q", id.RegistrationInfoName),
		fmt.Sprintf("Host Pool Name %q", id.HostPoolName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Host Pool Registration Info", segmentsStr)
}

func (id HostPoolRegistrationInfoId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DesktopVirtualization/hostPools/%s/registrationInfo/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.HostPoolName, id.RegistrationInfoName)
}

// HostPoolRegistrationInfoID parses a HostPoolRegistrationInfo ID into an HostPoolRegistrationInfoId struct
func HostPoolRegistrationInfoID(input string) (*HostPoolRegistrationInfoId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an HostPoolRegistrationInfo ID: %+v", input, err)
	}

	resourceId := HostPoolRegistrationInfoId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.HostPoolName, err = id.PopSegment("hostPools"); err != nil {
		return nil, err
	}
	if resourceId.RegistrationInfoName, err = id.PopSegment("registrationInfo"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// HostPoolRegistrationInfoIDInsensitively parses an HostPoolRegistrationInfo ID into an HostPoolRegistrationInfoId struct, insensitively
// This should only be used to parse an ID for rewriting, the HostPoolRegistrationInfoID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func HostPoolRegistrationInfoIDInsensitively(input string) (*HostPoolRegistrationInfoId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := HostPoolRegistrationInfoId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'hostPools' segment
	hostPoolsKey := "hostPools"
	for key := range id.Path {
		if strings.EqualFold(key, hostPoolsKey) {
			hostPoolsKey = key
			break
		}
	}
	if resourceId.HostPoolName, err = id.PopSegment(hostPoolsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'registrationInfo' segment
	registrationInfoKey := "registrationInfo"
	for key := range id.Path {
		if strings.EqualFold(key, registrationInfoKey) {
			registrationInfoKey = key
			break
		}
	}
	if resourceId.RegistrationInfoName, err = id.PopSegment(registrationInfoKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
