package commits

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&CommitId{})
}

var _ resourceids.ResourceId = &CommitId{}

// CommitId is a struct representing the Resource ID for a Commit
type CommitId struct {
	SubscriptionId     string
	ResourceGroupName  string
	NetworkManagerName string
	CommitName         string
}

// NewCommitID returns a new CommitId struct
func NewCommitID(subscriptionId string, resourceGroupName string, networkManagerName string, commitName string) CommitId {
	return CommitId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		NetworkManagerName: networkManagerName,
		CommitName:         commitName,
	}
}

// ParseCommitID parses 'input' into a CommitId
func ParseCommitID(input string) (*CommitId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CommitId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CommitId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCommitIDInsensitively parses 'input' case-insensitively into a CommitId
// note: this method should only be used for API response data and not user input
func ParseCommitIDInsensitively(input string) (*CommitId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CommitId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CommitId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CommitId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.NetworkManagerName, ok = input.Parsed["networkManagerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "networkManagerName", input)
	}

	if id.CommitName, ok = input.Parsed["commitName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "commitName", input)
	}

	return nil
}

// ValidateCommitID checks that 'input' can be parsed as a Commit ID
func ValidateCommitID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCommitID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Commit ID
func (id CommitId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkManagers/%s/commits/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName, id.CommitName)
}

// Segments returns a slice of Resource ID Segments which comprise this Commit ID
func (id CommitId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkManagers", "networkManagers", "networkManagers"),
		resourceids.UserSpecifiedSegment("networkManagerName", "networkManagerName"),
		resourceids.StaticSegment("staticCommits", "commits", "commits"),
		resourceids.UserSpecifiedSegment("commitName", "commitName"),
	}
}

// String returns a human-readable description of this Commit ID
func (id CommitId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Manager Name: %q", id.NetworkManagerName),
		fmt.Sprintf("Commit Name: %q", id.CommitName),
	}
	return fmt.Sprintf("Commit (%s)", strings.Join(components, "\n"))
}
