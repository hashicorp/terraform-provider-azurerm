package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SpringCloudAppBindingId struct {
	SubscriptionId string
	ResourceGroup  string
	SpringName     string
	AppName        string
	BindingName    string
}

func NewSpringCloudAppBindingID(subscriptionId, resourceGroup, springName, appName, bindingName string) SpringCloudAppBindingId {
	return SpringCloudAppBindingId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		SpringName:     springName,
		AppName:        appName,
		BindingName:    bindingName,
	}
}

func (id SpringCloudAppBindingId) String() string {
	segments := []string{
		fmt.Sprintf("Binding Name %q", id.BindingName),
		fmt.Sprintf("App Name %q", id.AppName),
		fmt.Sprintf("Spring Name %q", id.SpringName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Spring Cloud App Binding", segmentsStr)
}

func (id SpringCloudAppBindingId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/Spring/%s/apps/%s/bindings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SpringName, id.AppName, id.BindingName)
}

// SpringCloudAppBindingID parses a SpringCloudAppBinding ID into an SpringCloudAppBindingId struct
func SpringCloudAppBindingID(input string) (*SpringCloudAppBindingId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SpringCloudAppBindingId{
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
