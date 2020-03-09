package parse

import (
	"fmt"
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type BlueprintDefinitionId struct {
	ID   string
	Name string
	BlueprintDefinitionScopeId
}

func BlueprintDefinitionID(input string) (*BlueprintDefinitionId, error) {
	regex := regexp.MustCompile(`/providers/Microsoft\.Blueprint/blueprints/`)
	if !regex.MatchString(input) {
		return nil, fmt.Errorf("unable to parse Blueprint Definition ID %q", input)
	}

	segments := regex.Split(input, -1)

	if len(segments) != 2 {
		return nil, fmt.Errorf("unable to parse Blueprint Definition ID %q: Expected 2 segments after splition", input)
	}

	scope := segments[0]
	scopeId, err := BlueprintDefinitionScopeID(scope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Blueprint Definition ID %q: %+v", input, err)
	}

	name := segments[1]
	if name == "" {
		return nil, fmt.Errorf("unable to parse Blueprint Definition ID %q: blueprint name is empty", input)
	}

	id := BlueprintDefinitionId{
		ID:                         input,
		Name:                       name,
		BlueprintDefinitionScopeId: *scopeId,
	}

	return &id, nil
}

type ScopeType int

const (
	AtSubscription ScopeType = iota
	AtManagementGroup
)

type BlueprintDefinitionScopeId struct {
	ScopeId           string
	Type              ScopeType
	SubscriptionId    string
	ManagementGroupId string
}

func BlueprintDefinitionScopeID(input string) (*BlueprintDefinitionScopeId, error) {
	if input == "" {
		return nil, fmt.Errorf("unable to parse Blueprint Definition Scope ID: input is empty")
	}
	scopeId := BlueprintDefinitionScopeId{
		ScopeId: input,
	}

	if isManagementGroupId(input) {
		managementGroupId, _ := ManagementGroupID(input) // if this is a management group ID, there should not be any error.
		scopeId.ManagementGroupId = managementGroupId.GroupId
		scopeId.Type = AtManagementGroup
	} else {
		id, err := azure.ParseAzureResourceID(input)
		if err != nil {
			return nil, fmt.Errorf("unable to parse Blueprint Definition Scope ID %q: %+v", input, err)
		}
		if err := id.ValidateNoEmptySegments(input); err != nil {
			return nil, fmt.Errorf("unable to parse Blueprint Definition Scope ID %q: %+v", input, err)
		}
		if id.ResourceGroup != "" {
			return nil, fmt.Errorf("unable to parse Blueprint Definition Scope ID %q: scope cannot have resource groups", input)
		}
		scopeId.SubscriptionId = id.SubscriptionID
		scopeId.Type = AtSubscription
	}

	return &scopeId, nil
}

type PublishedBlueprintId struct {
	Version string
	BlueprintDefinitionId
}

func PublishedBlueprintID(input string) (*PublishedBlueprintId, error) {
	regex := regexp.MustCompile(`/versions/`)
	if !regex.MatchString(input) {
		return nil, fmt.Errorf("unable to parse Published Blueprint ID %q", input)
	}

	segments := regex.Split(input, -1)

	if len(segments) != 2 {
		return nil, fmt.Errorf("unable to parse Published Blueprint ID %q: Expected 2 segments after splition", input)
	}

	definitionId, err := BlueprintDefinitionID(segments[0])
	if err != nil {
		return nil, fmt.Errorf("unable to parse Published Blueprint ID %q: %+v", input, err)
	}

	version := segments[1]
	if version == "" {
		return nil, fmt.Errorf("unable to parse Published Blueprint ID %q: version is empty", input)
	}

	return &PublishedBlueprintId{
		Version:               version,
		BlueprintDefinitionId: *definitionId,
	}, nil
}

func isManagementGroupId(input string) bool {
	_, err := ManagementGroupID(input)
	return err == nil
}

// TODO -- move this to management group RP directory
type ManagementGroupId struct {
	GroupId string
}

func ManagementGroupID(input string) (*ManagementGroupId, error) {
	regex := regexp.MustCompile(`^/providers/[Mm]icrosoft\.[Mm]anagement/management[Gg]roups/`)
	if !regex.MatchString(input) {
		return nil, fmt.Errorf("unable to parse Management Group ID %q", input)
	}

	// Split the input ID by the regex
	segments := regex.Split(input, -1)
	if len(segments) != 2 {
		return nil, fmt.Errorf("unable to parse Management Group ID %q: expected id to have two segments after splitting", input)
	}

	groupID := segments[1]
	if groupID == "" {
		return nil, fmt.Errorf("unable to parse Management Group ID %q: group ID is empty", input)
	}

	id := ManagementGroupId{
		GroupId: groupID,
	}

	return &id, nil
}
