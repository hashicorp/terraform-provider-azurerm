// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AppServiceSlotCustomHostnameBindingId struct {
	SubscriptionId      string
	ResourceGroup       string
	SiteName            string
	SlotName            string
	HostNameBindingName string
}

func NewAppServiceSlotCustomHostnameBindingID(subscriptionId, resourceGroup, siteName, slotName, hostNameBindingName string) AppServiceSlotCustomHostnameBindingId {
	return AppServiceSlotCustomHostnameBindingId{
		SubscriptionId:      subscriptionId,
		ResourceGroup:       resourceGroup,
		SiteName:            siteName,
		SlotName:            slotName,
		HostNameBindingName: hostNameBindingName,
	}
}

func (id AppServiceSlotCustomHostnameBindingId) String() string {
	segments := []string{
		fmt.Sprintf("Host Name Binding Name %q", id.HostNameBindingName),
		fmt.Sprintf("Slot Name %q", id.SlotName),
		fmt.Sprintf("Site Name %q", id.SiteName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "App Service Slot Custom Hostname Binding", segmentsStr)
}

func (id AppServiceSlotCustomHostnameBindingId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/slots/%s/hostNameBindings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SiteName, id.SlotName, id.HostNameBindingName)
}

// AppServiceSlotCustomHostnameBindingID parses a AppServiceSlotCustomHostnameBinding ID into an AppServiceSlotCustomHostnameBindingId struct
func AppServiceSlotCustomHostnameBindingID(input string) (*AppServiceSlotCustomHostnameBindingId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an AppServiceSlotCustomHostnameBinding ID: %+v", input, err)
	}

	resourceId := AppServiceSlotCustomHostnameBindingId{
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
	if resourceId.SlotName, err = id.PopSegment("slots"); err != nil {
		return nil, err
	}
	if resourceId.HostNameBindingName, err = id.PopSegment("hostNameBindings"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
