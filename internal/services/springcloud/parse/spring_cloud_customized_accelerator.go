// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SpringCloudCustomizedAcceleratorId struct {
	SubscriptionId             string
	ResourceGroup              string
	SpringName                 string
	ApplicationAcceleratorName string
	CustomizedAcceleratorName  string
}

func NewSpringCloudCustomizedAcceleratorID(subscriptionId, resourceGroup, springName, applicationAcceleratorName, customizedAcceleratorName string) SpringCloudCustomizedAcceleratorId {
	return SpringCloudCustomizedAcceleratorId{
		SubscriptionId:             subscriptionId,
		ResourceGroup:              resourceGroup,
		SpringName:                 springName,
		ApplicationAcceleratorName: applicationAcceleratorName,
		CustomizedAcceleratorName:  customizedAcceleratorName,
	}
}

func (id SpringCloudCustomizedAcceleratorId) String() string {
	segments := []string{
		fmt.Sprintf("Customized Accelerator Name %q", id.CustomizedAcceleratorName),
		fmt.Sprintf("Application Accelerator Name %q", id.ApplicationAcceleratorName),
		fmt.Sprintf("Spring Name %q", id.SpringName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Spring Cloud Customized Accelerator", segmentsStr)
}

func (id SpringCloudCustomizedAcceleratorId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/Spring/%s/applicationAccelerators/%s/customizedAccelerators/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SpringName, id.ApplicationAcceleratorName, id.CustomizedAcceleratorName)
}

// SpringCloudCustomizedAcceleratorID parses a SpringCloudCustomizedAccelerator ID into an SpringCloudCustomizedAcceleratorId struct
func SpringCloudCustomizedAcceleratorID(input string) (*SpringCloudCustomizedAcceleratorId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SpringCloudCustomizedAccelerator ID: %+v", input, err)
	}

	resourceId := SpringCloudCustomizedAcceleratorId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.SpringName, err = id.PopSegment("Spring"); err != nil {
		return nil, err
	}
	if resourceId.ApplicationAcceleratorName, err = id.PopSegment("applicationAccelerators"); err != nil {
		return nil, err
	}
	if resourceId.CustomizedAcceleratorName, err = id.PopSegment("customizedAccelerators"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
