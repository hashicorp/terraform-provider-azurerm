package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SpringCloudApplicationPerformanceMonitoringId struct {
	SubscriptionId string
	ResourceGroup  string
	SpringName     string
	ApmName        string
}

func NewSpringCloudApplicationPerformanceMonitoringID(subscriptionId, resourceGroup, springName, apmName string) SpringCloudApplicationPerformanceMonitoringId {
	return SpringCloudApplicationPerformanceMonitoringId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		SpringName:     springName,
		ApmName:        apmName,
	}
}

func (id SpringCloudApplicationPerformanceMonitoringId) String() string {
	segments := []string{
		fmt.Sprintf("Apm Name %q", id.ApmName),
		fmt.Sprintf("Spring Name %q", id.SpringName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Spring Cloud Application Performance Monitoring", segmentsStr)
}

func (id SpringCloudApplicationPerformanceMonitoringId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/apms/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SpringName, id.ApmName)
}

// SpringCloudApplicationPerformanceMonitoringID parses a SpringCloudApplicationPerformanceMonitoring ID into an SpringCloudApplicationPerformanceMonitoringId struct
func SpringCloudApplicationPerformanceMonitoringID(input string) (*SpringCloudApplicationPerformanceMonitoringId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SpringCloudApplicationPerformanceMonitoring ID: %+v", input, err)
	}

	resourceId := SpringCloudApplicationPerformanceMonitoringId{
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
	if resourceId.ApmName, err = id.PopSegment("apms"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// SpringCloudApplicationPerformanceMonitoringIDInsensitively parses an SpringCloudApplicationPerformanceMonitoring ID into an SpringCloudApplicationPerformanceMonitoringId struct, insensitively
// This should only be used to parse an ID for rewriting, the SpringCloudApplicationPerformanceMonitoringID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func SpringCloudApplicationPerformanceMonitoringIDInsensitively(input string) (*SpringCloudApplicationPerformanceMonitoringId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SpringCloudApplicationPerformanceMonitoringId{
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

	// find the correct casing for the 'apms' segment
	apmsKey := "apms"
	for key := range id.Path {
		if strings.EqualFold(key, apmsKey) {
			apmsKey = key
			break
		}
	}
	if resourceId.ApmName, err = id.PopSegment(apmsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
