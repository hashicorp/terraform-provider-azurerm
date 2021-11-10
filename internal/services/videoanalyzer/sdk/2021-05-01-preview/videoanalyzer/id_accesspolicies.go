package videoanalyzer

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AccessPoliciesId struct {
	SubscriptionId    string
	ResourceGroup     string
	VideoAnalyzerName string
	AccessPolicyName  string
}

func NewAccessPoliciesID(subscriptionId, resourceGroup, videoAnalyzerName, accessPolicyName string) AccessPoliciesId {
	return AccessPoliciesId{
		SubscriptionId:    subscriptionId,
		ResourceGroup:     resourceGroup,
		VideoAnalyzerName: videoAnalyzerName,
		AccessPolicyName:  accessPolicyName,
	}
}

func (id AccessPoliciesId) String() string {
	segments := []string{
		fmt.Sprintf("Access Policy Name %q", id.AccessPolicyName),
		fmt.Sprintf("Video Analyzer Name %q", id.VideoAnalyzerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Access Policies", segmentsStr)
}

func (id AccessPoliciesId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Media/videoAnalyzers/%s/accessPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VideoAnalyzerName, id.AccessPolicyName)
}

// ParseAccessPoliciesID parses a AccessPolicies ID into an AccessPoliciesId struct
func ParseAccessPoliciesID(input string) (*AccessPoliciesId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AccessPoliciesId{
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
	if resourceId.AccessPolicyName, err = id.PopSegment("accessPolicies"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParseAccessPoliciesIDInsensitively parses an AccessPolicies ID into an AccessPoliciesId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseAccessPoliciesID method should be used instead for validation etc.
func ParseAccessPoliciesIDInsensitively(input string) (*AccessPoliciesId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AccessPoliciesId{
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

	// find the correct casing for the 'accessPolicies' segment
	accessPoliciesKey := "accessPolicies"
	for key := range id.Path {
		if strings.EqualFold(key, accessPoliciesKey) {
			accessPoliciesKey = key
			break
		}
	}
	if resourceId.AccessPolicyName, err = id.PopSegment(accessPoliciesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
