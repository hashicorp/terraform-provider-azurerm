package longrunningbackup

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&BackupsV2Id{})
}

var _ resourceids.ResourceId = &BackupsV2Id{}

// BackupsV2Id is a struct representing the Resource ID for a Backups V 2
type BackupsV2Id struct {
	SubscriptionId     string
	ResourceGroupName  string
	FlexibleServerName string
	BackupsV2Name      string
}

// NewBackupsV2ID returns a new BackupsV2Id struct
func NewBackupsV2ID(subscriptionId string, resourceGroupName string, flexibleServerName string, backupsV2Name string) BackupsV2Id {
	return BackupsV2Id{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		FlexibleServerName: flexibleServerName,
		BackupsV2Name:      backupsV2Name,
	}
}

// ParseBackupsV2ID parses 'input' into a BackupsV2Id
func ParseBackupsV2ID(input string) (*BackupsV2Id, error) {
	parser := resourceids.NewParserFromResourceIdType(&BackupsV2Id{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BackupsV2Id{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseBackupsV2IDInsensitively parses 'input' case-insensitively into a BackupsV2Id
// note: this method should only be used for API response data and not user input
func ParseBackupsV2IDInsensitively(input string) (*BackupsV2Id, error) {
	parser := resourceids.NewParserFromResourceIdType(&BackupsV2Id{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BackupsV2Id{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *BackupsV2Id) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.FlexibleServerName, ok = input.Parsed["flexibleServerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "flexibleServerName", input)
	}

	if id.BackupsV2Name, ok = input.Parsed["backupsV2Name"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "backupsV2Name", input)
	}

	return nil
}

// ValidateBackupsV2ID checks that 'input' can be parsed as a Backups V 2 ID
func ValidateBackupsV2ID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBackupsV2ID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Backups V 2 ID
func (id BackupsV2Id) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DBforMySQL/flexibleServers/%s/backupsV2/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.FlexibleServerName, id.BackupsV2Name)
}

// Segments returns a slice of Resource ID Segments which comprise this Backups V 2 ID
func (id BackupsV2Id) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDBforMySQL", "Microsoft.DBforMySQL", "Microsoft.DBforMySQL"),
		resourceids.StaticSegment("staticFlexibleServers", "flexibleServers", "flexibleServers"),
		resourceids.UserSpecifiedSegment("flexibleServerName", "flexibleServerName"),
		resourceids.StaticSegment("staticBackupsV2", "backupsV2", "backupsV2"),
		resourceids.UserSpecifiedSegment("backupsV2Name", "backupsV2Name"),
	}
}

// String returns a human-readable description of this Backups V 2 ID
func (id BackupsV2Id) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Flexible Server Name: %q", id.FlexibleServerName),
		fmt.Sprintf("Backups V 2 Name: %q", id.BackupsV2Name),
	}
	return fmt.Sprintf("Backups V 2 (%s)", strings.Join(components, "\n"))
}
