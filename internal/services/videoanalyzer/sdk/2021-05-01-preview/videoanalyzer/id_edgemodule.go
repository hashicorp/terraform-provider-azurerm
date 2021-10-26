package videoanalyzer

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type EdgeModuleId struct {
	SubscriptionId    string
	ResourceGroup     string
	VideoAnalyzerName string
	Name              string
}

func NewEdgeModuleID(subscriptionId, resourceGroup, videoAnalyzerName, name string) EdgeModuleId {
	return EdgeModuleId{
		SubscriptionId:    subscriptionId,
		ResourceGroup:     resourceGroup,
		VideoAnalyzerName: videoAnalyzerName,
		Name:              name,
	}
}

func (id EdgeModuleId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Video Analyzer Name %q", id.VideoAnalyzerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Edge Module", segmentsStr)
}

func (id EdgeModuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Media/videoAnalyzers/%s/edgeModules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VideoAnalyzerName, id.Name)
}

// ParseEdgeModuleID parses a EdgeModule ID into an EdgeModuleId struct
func ParseEdgeModuleID(input string) (*EdgeModuleId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := EdgeModuleId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.VideoAnalyzerName, err = id.PopSegment("videoAnalyzers"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("edgeModules"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParseEdgeModuleIDInsensitively parses an EdgeModule ID into an EdgeModuleId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseEdgeModuleID method should be used instead for validation etc.
func ParseEdgeModuleIDInsensitively(input string) (*EdgeModuleId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := EdgeModuleId{
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
	if resourceId.VideoAnalyzerName, err = id.PopSegment(videoAnalyzersKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'edgeModules' segment
	edgeModulesKey := "edgeModules"
	for key := range id.Path {
		if strings.EqualFold(key, edgeModulesKey) {
			edgeModulesKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(edgeModulesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
