package tenantconfiguration

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ConfigurationId{}

// ConfigurationId is a struct representing the Resource ID for a Configuration
type ConfigurationId struct {
	ConfigurationName ConfigurationName
}

// NewConfigurationID returns a new ConfigurationId struct
func NewConfigurationID(configurationName ConfigurationName) ConfigurationId {
	return ConfigurationId{
		ConfigurationName: configurationName,
	}
}

// ParseConfigurationID parses 'input' into a ConfigurationId
func ParseConfigurationID(input string) (*ConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(ConfigurationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConfigurationId{}

	if v, ok := parsed.Parsed["configurationName"]; true {
		if !ok {
			return nil, fmt.Errorf("the segment 'configurationName' was not found in the resource id %q", input)
		}

		configurationName, err := parseConfigurationName(v)
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %+v", v, err)
		}
		id.ConfigurationName = *configurationName
	}

	return &id, nil
}

// ParseConfigurationIDInsensitively parses 'input' case-insensitively into a ConfigurationId
// note: this method should only be used for API response data and not user input
func ParseConfigurationIDInsensitively(input string) (*ConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(ConfigurationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConfigurationId{}

	if v, ok := parsed.Parsed["configurationName"]; true {
		if !ok {
			return nil, fmt.Errorf("the segment 'configurationName' was not found in the resource id %q", input)
		}

		configurationName, err := parseConfigurationName(v)
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %+v", v, err)
		}
		id.ConfigurationName = *configurationName
	}

	return &id, nil
}

// ValidateConfigurationID checks that 'input' can be parsed as a Configuration ID
func ValidateConfigurationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseConfigurationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Configuration ID
func (id ConfigurationId) ID() string {
	fmtString := "/providers/Microsoft.Portal/tenantConfigurations/%s"
	return fmt.Sprintf(fmtString, string(id.ConfigurationName))
}

// Segments returns a slice of Resource ID Segments which comprise this Configuration ID
func (id ConfigurationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftPortal", "Microsoft.Portal", "Microsoft.Portal"),
		resourceids.StaticSegment("staticTenantConfigurations", "tenantConfigurations", "tenantConfigurations"),
		resourceids.ConstantSegment("configurationName", PossibleValuesForConfigurationName(), "default"),
	}
}

// String returns a human-readable description of this Configuration ID
func (id ConfigurationId) String() string {
	components := []string{
		fmt.Sprintf("Configuration Name: %q", string(id.ConfigurationName)),
	}
	return fmt.Sprintf("Configuration (%s)", strings.Join(components, "\n"))
}
