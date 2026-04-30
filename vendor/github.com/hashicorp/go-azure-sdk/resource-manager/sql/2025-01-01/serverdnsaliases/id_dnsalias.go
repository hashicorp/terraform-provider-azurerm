package serverdnsaliases

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DnsAliasId{})
}

var _ resourceids.ResourceId = &DnsAliasId{}

// DnsAliasId is a struct representing the Resource ID for a Dns Alias
type DnsAliasId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServerName        string
	DnsAliasName      string
}

// NewDnsAliasID returns a new DnsAliasId struct
func NewDnsAliasID(subscriptionId string, resourceGroupName string, serverName string, dnsAliasName string) DnsAliasId {
	return DnsAliasId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServerName:        serverName,
		DnsAliasName:      dnsAliasName,
	}
}

// ParseDnsAliasID parses 'input' into a DnsAliasId
func ParseDnsAliasID(input string) (*DnsAliasId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DnsAliasId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DnsAliasId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDnsAliasIDInsensitively parses 'input' case-insensitively into a DnsAliasId
// note: this method should only be used for API response data and not user input
func ParseDnsAliasIDInsensitively(input string) (*DnsAliasId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DnsAliasId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DnsAliasId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DnsAliasId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ServerName, ok = input.Parsed["serverName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "serverName", input)
	}

	if id.DnsAliasName, ok = input.Parsed["dnsAliasName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "dnsAliasName", input)
	}

	return nil
}

// ValidateDnsAliasID checks that 'input' can be parsed as a Dns Alias ID
func ValidateDnsAliasID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDnsAliasID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Dns Alias ID
func (id DnsAliasId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/servers/%s/dnsAliases/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServerName, id.DnsAliasName)
}

// Segments returns a slice of Resource ID Segments which comprise this Dns Alias ID
func (id DnsAliasId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSql", "Microsoft.Sql", "Microsoft.Sql"),
		resourceids.StaticSegment("staticServers", "servers", "servers"),
		resourceids.UserSpecifiedSegment("serverName", "serverName"),
		resourceids.StaticSegment("staticDnsAliases", "dnsAliases", "dnsAliases"),
		resourceids.UserSpecifiedSegment("dnsAliasName", "dnsAliasName"),
	}
}

// String returns a human-readable description of this Dns Alias ID
func (id DnsAliasId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Server Name: %q", id.ServerName),
		fmt.Sprintf("Dns Alias Name: %q", id.DnsAliasName),
	}
	return fmt.Sprintf("Dns Alias (%s)", strings.Join(components, "\n"))
}
