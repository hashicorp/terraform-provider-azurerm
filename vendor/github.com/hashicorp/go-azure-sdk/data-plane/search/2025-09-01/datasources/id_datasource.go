package datasources

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &DatasourceId{}

// DatasourceId is a struct representing the Resource ID for a Datasource
type DatasourceId struct {
	BaseURI        string
	DatasourceName string
}

// NewDatasourceID returns a new DatasourceId struct
func NewDatasourceID(baseURI string, datasourceName string) DatasourceId {
	return DatasourceId{
		BaseURI:        strings.TrimSuffix(baseURI, "/"),
		DatasourceName: datasourceName,
	}
}

// ParseDatasourceID parses 'input' into a DatasourceId
func ParseDatasourceID(input string) (*DatasourceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DatasourceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DatasourceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDatasourceIDInsensitively parses 'input' case-insensitively into a DatasourceId
// note: this method should only be used for API response data and not user input
func ParseDatasourceIDInsensitively(input string) (*DatasourceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DatasourceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DatasourceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DatasourceId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BaseURI, ok = input.Parsed["baseURI"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "baseURI", input)
	}

	if id.DatasourceName, ok = input.Parsed["datasourceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "datasourceName", input)
	}

	return nil
}

// ValidateDatasourceID checks that 'input' can be parsed as a Datasource ID
func ValidateDatasourceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDatasourceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Datasource ID
func (id DatasourceId) ID() string {
	fmtString := "%s/datasources/%s"
	return fmt.Sprintf(fmtString, strings.TrimSuffix(id.BaseURI, "/"), id.DatasourceName)
}

// Path returns the formatted Datasource ID without the BaseURI
func (id DatasourceId) Path() string {
	fmtString := "/datasources/%s"
	return fmt.Sprintf(fmtString, id.DatasourceName)
}

// PathElements returns the values of Datasource ID Segments without the BaseURI
func (id DatasourceId) PathElements() []any {
	return []any{id.DatasourceName}
}

// Segments returns a slice of Resource ID Segments which comprise this Datasource ID
func (id DatasourceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.DataPlaneBaseURISegment("baseURI", "https://endpoint-url.example.com"),
		resourceids.StaticSegment("datasources", "datasources", "datasources"),
		resourceids.UserSpecifiedSegment("datasourceName", "datasourceName"),
	}
}

// String returns a human-readable description of this Datasource ID
func (id DatasourceId) String() string {
	components := []string{
		fmt.Sprintf("Base URI: %q", id.BaseURI),
		fmt.Sprintf("Datasource Name: %q", id.DatasourceName),
	}
	return fmt.Sprintf("Datasource (%s)", strings.Join(components, "\n"))
}
