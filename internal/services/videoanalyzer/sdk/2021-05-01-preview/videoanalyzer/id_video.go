package videoanalyzer

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = VideoId{}

// VideoId is a struct representing the Resource ID for a Video
type VideoId struct {
	SubscriptionId    string
	ResourceGroupName string
	AccountName       string
	VideoName         string
}

// NewVideoID returns a new VideoId struct
func NewVideoID(subscriptionId string, resourceGroupName string, accountName string, videoName string) VideoId {
	return VideoId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		AccountName:       accountName,
		VideoName:         videoName,
	}
}

// ParseVideoID parses 'input' into a VideoId
func ParseVideoID(input string) (*VideoId, error) {
	parser := resourceids.NewParserFromResourceIdType(VideoId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VideoId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.VideoName, ok = parsed.Parsed["videoName"]; !ok {
		return nil, fmt.Errorf("the segment 'videoName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseVideoIDInsensitively parses 'input' case-insensitively into a VideoId
// note: this method should only be used for API response data and not user input
func ParseVideoIDInsensitively(input string) (*VideoId, error) {
	parser := resourceids.NewParserFromResourceIdType(VideoId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VideoId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.VideoName, ok = parsed.Parsed["videoName"]; !ok {
		return nil, fmt.Errorf("the segment 'videoName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateVideoID checks that 'input' can be parsed as a Video ID
func ValidateVideoID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVideoID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Video ID
func (id VideoId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Media/videoAnalyzers/%s/videos/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.VideoName)
}

// Segments returns a slice of Resource ID Segments which comprise this Video ID
func (id VideoId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMedia", "Microsoft.Media", "Microsoft.Media"),
		resourceids.StaticSegment("staticVideoAnalyzers", "videoAnalyzers", "videoAnalyzers"),
		resourceids.UserSpecifiedSegment("accountName", "accountValue"),
		resourceids.StaticSegment("staticVideos", "videos", "videos"),
		resourceids.UserSpecifiedSegment("videoName", "videoValue"),
	}
}

// String returns a human-readable description of this Video ID
func (id VideoId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Video Name: %q", id.VideoName),
	}
	return fmt.Sprintf("Video (%s)", strings.Join(components, "\n"))
}
