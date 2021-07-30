package clusters

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ClusterId struct {
	SubscriptionId   string
	ResourceGroup    string
	PrivateCloudName string
	Name             string
}

func NewClusterID(subscriptionId, resourceGroup, privateCloudName, name string) ClusterId {
	return ClusterId{
		SubscriptionId:   subscriptionId,
		ResourceGroup:    resourceGroup,
		PrivateCloudName: privateCloudName,
		Name:             name,
	}
}

func (id ClusterId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Private Cloud Name %q", id.PrivateCloudName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Cluster", segmentsStr)
}

func (id ClusterId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AVS/privateClouds/%s/clusters/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.PrivateCloudName, id.Name)
}

// ParseClusterID parses a Cluster ID into an ClusterId struct
func ParseClusterID(input string) (*ClusterId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ClusterId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.PrivateCloudName, err = id.PopSegment("privateClouds"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("clusters"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParseClusterIDInsensitively parses an Cluster ID into an ClusterId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseClusterID method should be used instead for validation etc.
func ParseClusterIDInsensitively(input string) (*ClusterId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ClusterId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'privateClouds' segment
	privateCloudsKey := "privateClouds"
	for key := range id.Path {
		if strings.EqualFold(key, privateCloudsKey) {
			privateCloudsKey = key
			break
		}
	}
	if resourceId.PrivateCloudName, err = id.PopSegment(privateCloudsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'clusters' segment
	clustersKey := "clusters"
	for key := range id.Path {
		if strings.EqualFold(key, clustersKey) {
			clustersKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(clustersKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
