package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type JobId struct {
	SubscriptionId   string
	ResourceGroup    string
	MediaserviceName string
	TransformName    string
	Name             string
}

func NewJobID(subscriptionId, resourceGroup, mediaserviceName, transformName, name string) JobId {
	return JobId{
		SubscriptionId:   subscriptionId,
		ResourceGroup:    resourceGroup,
		MediaserviceName: mediaserviceName,
		TransformName:    transformName,
		Name:             name,
	}
}

func (id JobId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Transform Name %q", id.TransformName),
		fmt.Sprintf("Mediaservice Name %q", id.MediaserviceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Job", segmentsStr)
}

func (id JobId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Media/mediaservices/%s/transforms/%s/jobs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.MediaserviceName, id.TransformName, id.Name)
}

// JobID parses a Job ID into an JobId struct
func JobID(input string) (*JobId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := JobId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.MediaserviceName, err = id.PopSegment("mediaservices"); err != nil {
		return nil, err
	}
	if resourceId.TransformName, err = id.PopSegment("transforms"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("jobs"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
