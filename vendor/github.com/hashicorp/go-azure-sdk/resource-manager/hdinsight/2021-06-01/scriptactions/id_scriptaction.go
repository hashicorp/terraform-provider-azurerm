package scriptactions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ScriptActionId{})
}

var _ resourceids.ResourceId = &ScriptActionId{}

// ScriptActionId is a struct representing the Resource ID for a Script Action
type ScriptActionId struct {
	SubscriptionId    string
	ResourceGroupName string
	ClusterName       string
	ScriptActionName  string
}

// NewScriptActionID returns a new ScriptActionId struct
func NewScriptActionID(subscriptionId string, resourceGroupName string, clusterName string, scriptActionName string) ScriptActionId {
	return ScriptActionId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ClusterName:       clusterName,
		ScriptActionName:  scriptActionName,
	}
}

// ParseScriptActionID parses 'input' into a ScriptActionId
func ParseScriptActionID(input string) (*ScriptActionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScriptActionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScriptActionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseScriptActionIDInsensitively parses 'input' case-insensitively into a ScriptActionId
// note: this method should only be used for API response data and not user input
func ParseScriptActionIDInsensitively(input string) (*ScriptActionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScriptActionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScriptActionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ScriptActionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ClusterName, ok = input.Parsed["clusterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "clusterName", input)
	}

	if id.ScriptActionName, ok = input.Parsed["scriptActionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scriptActionName", input)
	}

	return nil
}

// ValidateScriptActionID checks that 'input' can be parsed as a Script Action ID
func ValidateScriptActionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScriptActionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Script Action ID
func (id ScriptActionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.HDInsight/clusters/%s/scriptActions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ClusterName, id.ScriptActionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Script Action ID
func (id ScriptActionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftHDInsight", "Microsoft.HDInsight", "Microsoft.HDInsight"),
		resourceids.StaticSegment("staticClusters", "clusters", "clusters"),
		resourceids.UserSpecifiedSegment("clusterName", "clusterName"),
		resourceids.StaticSegment("staticScriptActions", "scriptActions", "scriptActions"),
		resourceids.UserSpecifiedSegment("scriptActionName", "scriptActionName"),
	}
}

// String returns a human-readable description of this Script Action ID
func (id ScriptActionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cluster Name: %q", id.ClusterName),
		fmt.Sprintf("Script Action Name: %q", id.ScriptActionName),
	}
	return fmt.Sprintf("Script Action (%s)", strings.Join(components, "\n"))
}
