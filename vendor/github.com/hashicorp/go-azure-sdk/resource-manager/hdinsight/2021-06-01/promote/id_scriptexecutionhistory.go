package promote

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ScriptExecutionHistoryId{})
}

var _ resourceids.ResourceId = &ScriptExecutionHistoryId{}

// ScriptExecutionHistoryId is a struct representing the Resource ID for a Script Execution History
type ScriptExecutionHistoryId struct {
	SubscriptionId    string
	ResourceGroupName string
	ClusterName       string
	ScriptExecutionId string
}

// NewScriptExecutionHistoryID returns a new ScriptExecutionHistoryId struct
func NewScriptExecutionHistoryID(subscriptionId string, resourceGroupName string, clusterName string, scriptExecutionId string) ScriptExecutionHistoryId {
	return ScriptExecutionHistoryId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ClusterName:       clusterName,
		ScriptExecutionId: scriptExecutionId,
	}
}

// ParseScriptExecutionHistoryID parses 'input' into a ScriptExecutionHistoryId
func ParseScriptExecutionHistoryID(input string) (*ScriptExecutionHistoryId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScriptExecutionHistoryId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScriptExecutionHistoryId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseScriptExecutionHistoryIDInsensitively parses 'input' case-insensitively into a ScriptExecutionHistoryId
// note: this method should only be used for API response data and not user input
func ParseScriptExecutionHistoryIDInsensitively(input string) (*ScriptExecutionHistoryId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScriptExecutionHistoryId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScriptExecutionHistoryId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ScriptExecutionHistoryId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.ScriptExecutionId, ok = input.Parsed["scriptExecutionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scriptExecutionId", input)
	}

	return nil
}

// ValidateScriptExecutionHistoryID checks that 'input' can be parsed as a Script Execution History ID
func ValidateScriptExecutionHistoryID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScriptExecutionHistoryID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Script Execution History ID
func (id ScriptExecutionHistoryId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.HDInsight/clusters/%s/scriptExecutionHistory/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ClusterName, id.ScriptExecutionId)
}

// Segments returns a slice of Resource ID Segments which comprise this Script Execution History ID
func (id ScriptExecutionHistoryId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftHDInsight", "Microsoft.HDInsight", "Microsoft.HDInsight"),
		resourceids.StaticSegment("staticClusters", "clusters", "clusters"),
		resourceids.UserSpecifiedSegment("clusterName", "clusterName"),
		resourceids.StaticSegment("staticScriptExecutionHistory", "scriptExecutionHistory", "scriptExecutionHistory"),
		resourceids.UserSpecifiedSegment("scriptExecutionId", "scriptExecutionId"),
	}
}

// String returns a human-readable description of this Script Execution History ID
func (id ScriptExecutionHistoryId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cluster Name: %q", id.ClusterName),
		fmt.Sprintf("Script Execution: %q", id.ScriptExecutionId),
	}
	return fmt.Sprintf("Script Execution History (%s)", strings.Join(components, "\n"))
}
