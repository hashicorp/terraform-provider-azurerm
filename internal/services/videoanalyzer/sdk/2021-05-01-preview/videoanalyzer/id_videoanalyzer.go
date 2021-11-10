package videoanalyzer

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type VideoAnalyzerId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewVideoAnalyzerID(subscriptionId, resourceGroup, name string) VideoAnalyzerId {
	return VideoAnalyzerId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id VideoAnalyzerId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Video Analyzer", segmentsStr)
}

func (id VideoAnalyzerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Media/videoAnalyzers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

// ParseVideoAnalyzerID parses a VideoAnalyzer ID into an VideoAnalyzerId struct
func ParseVideoAnalyzerID(input string) (*VideoAnalyzerId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := VideoAnalyzerId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.Name, err = id.PopSegment("videoAnalyzers"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParseVideoAnalyzerIDInsensitively parses an VideoAnalyzer ID into an VideoAnalyzerId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseVideoAnalyzerID method should be used instead for validation etc.
func ParseVideoAnalyzerIDInsensitively(input string) (*VideoAnalyzerId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := VideoAnalyzerId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'videoAnalyzers' segment
	videoAnalyzersKey := "videoAnalyzers"
	for key := range id.Path {
		if strings.EqualFold(key, videoAnalyzersKey) {
			videoAnalyzersKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(videoAnalyzersKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
