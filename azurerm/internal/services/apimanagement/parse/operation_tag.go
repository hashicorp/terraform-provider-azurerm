package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type OperationTagId struct {
	SubscriptionId string
	ResourceGroup  string
	ServiceName    string
	ApiName        string
	OperationName  string
	TagName        string
}

func NewOperationTagID(subscriptionId, resourceGroup, serviceName, apiName, operationName, tagName string) OperationTagId {
	return OperationTagId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ServiceName:    serviceName,
		ApiName:        apiName,
		OperationName:  operationName,
		TagName:        tagName,
	}
}

func (id OperationTagId) String() string {
	segments := []string{
		fmt.Sprintf("Tag Name %q", id.TagName),
		fmt.Sprintf("Operation Name %q", id.OperationName),
		fmt.Sprintf("Api Name %q", id.ApiName),
		fmt.Sprintf("Service Name %q", id.ServiceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Operation Tag", segmentsStr)
}

func (id OperationTagId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/apis/%s/operations/%s/tags/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServiceName, id.ApiName, id.OperationName, id.TagName)
}

// OperationTagID parses a OperationTag ID into an OperationTagId struct
func OperationTagID(input string) (*OperationTagId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := OperationTagId{
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
	if resourceId.OperationName, err = id.PopSegment("operations"); err != nil {
		return nil, err
	}
	if resourceId.TagName, err = id.PopSegment("tags"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
