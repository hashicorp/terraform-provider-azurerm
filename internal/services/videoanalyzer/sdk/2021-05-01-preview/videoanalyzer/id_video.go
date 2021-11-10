package videoanalyzer

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type VideoId struct {
	SubscriptionId    string
	ResourceGroup     string
	VideoAnalyzerName string
	Name              string
}

func NewVideoID(subscriptionId, resourceGroup, videoAnalyzerName, name string) VideoId {
	return VideoId{
		SubscriptionId:    subscriptionId,
		ResourceGroup:     resourceGroup,
		VideoAnalyzerName: videoAnalyzerName,
		Name:              name,
	}
}

func (id VideoId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Video Analyzer Name %q", id.VideoAnalyzerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Video", segmentsStr)
}

func (id VideoId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Media/videoAnalyzers/%s/videos/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VideoAnalyzerName, id.Name)
}

// ParseVideoID parses a Video ID into an VideoId struct
func ParseVideoID(input string) (*VideoId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := VideoId{
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
	if resourceId.Name, err = id.PopSegment("videos"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParseVideoIDInsensitively parses an Video ID into an VideoId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseVideoID method should be used instead for validation etc.
func ParseVideoIDInsensitively(input string) (*VideoId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := VideoId{
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

	// find the correct casing for the 'videos' segment
	videosKey := "videos"
	for key := range id.Path {
		if strings.EqualFold(key, videosKey) {
			videosKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(videosKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
