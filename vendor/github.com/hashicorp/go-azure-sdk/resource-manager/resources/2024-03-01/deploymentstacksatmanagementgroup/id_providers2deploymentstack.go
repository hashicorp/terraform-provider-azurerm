package deploymentstacksatmanagementgroup

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&Providers2DeploymentStackId{})
}

var _ resourceids.ResourceId = &Providers2DeploymentStackId{}

// Providers2DeploymentStackId is a struct representing the Resource ID for a Providers 2 Deployment Stack
type Providers2DeploymentStackId struct {
	ManagementGroupId   string
	DeploymentStackName string
}

// NewProviders2DeploymentStackID returns a new Providers2DeploymentStackId struct
func NewProviders2DeploymentStackID(managementGroupId string, deploymentStackName string) Providers2DeploymentStackId {
	return Providers2DeploymentStackId{
		ManagementGroupId:   managementGroupId,
		DeploymentStackName: deploymentStackName,
	}
}

// ParseProviders2DeploymentStackID parses 'input' into a Providers2DeploymentStackId
func ParseProviders2DeploymentStackID(input string) (*Providers2DeploymentStackId, error) {
	parser := resourceids.NewParserFromResourceIdType(&Providers2DeploymentStackId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := Providers2DeploymentStackId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseProviders2DeploymentStackIDInsensitively parses 'input' case-insensitively into a Providers2DeploymentStackId
// note: this method should only be used for API response data and not user input
func ParseProviders2DeploymentStackIDInsensitively(input string) (*Providers2DeploymentStackId, error) {
	parser := resourceids.NewParserFromResourceIdType(&Providers2DeploymentStackId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := Providers2DeploymentStackId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *Providers2DeploymentStackId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.ManagementGroupId, ok = input.Parsed["managementGroupId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "managementGroupId", input)
	}

	if id.DeploymentStackName, ok = input.Parsed["deploymentStackName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "deploymentStackName", input)
	}

	return nil
}

// ValidateProviders2DeploymentStackID checks that 'input' can be parsed as a Providers 2 Deployment Stack ID
func ValidateProviders2DeploymentStackID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProviders2DeploymentStackID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Providers 2 Deployment Stack ID
func (id Providers2DeploymentStackId) ID() string {
	fmtString := "/providers/Microsoft.Management/managementGroups/%s/providers/Microsoft.Resources/deploymentStacks/%s"
	return fmt.Sprintf(fmtString, id.ManagementGroupId, id.DeploymentStackName)
}

// Segments returns a slice of Resource ID Segments which comprise this Providers 2 Deployment Stack ID
func (id Providers2DeploymentStackId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftManagement", "Microsoft.Management", "Microsoft.Management"),
		resourceids.StaticSegment("staticManagementGroups", "managementGroups", "managementGroups"),
		resourceids.UserSpecifiedSegment("managementGroupId", "managementGroupId"),
		resourceids.StaticSegment("staticProviders2", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftResources", "Microsoft.Resources", "Microsoft.Resources"),
		resourceids.StaticSegment("staticDeploymentStacks", "deploymentStacks", "deploymentStacks"),
		resourceids.UserSpecifiedSegment("deploymentStackName", "deploymentStackName"),
	}
}

// String returns a human-readable description of this Providers 2 Deployment Stack ID
func (id Providers2DeploymentStackId) String() string {
	components := []string{
		fmt.Sprintf("Management Group: %q", id.ManagementGroupId),
		fmt.Sprintf("Deployment Stack Name: %q", id.DeploymentStackName),
	}
	return fmt.Sprintf("Providers 2 Deployment Stack (%s)", strings.Join(components, "\n"))
}
