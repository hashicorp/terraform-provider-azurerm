package pluginsdk

import (
	"fmt"
	"slices"
	"strings"
	"unicode"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func convertToSnakeCase(input string) string {
	w := []rune(input)
	output := ""
	for _, r := range w {
		if unicode.IsUpper(r) {
			output += "_"
		}
		output += string(r)
	}

	return strings.ToLower(output)
}

func segmentTypeSupported(segment resourceids.SegmentType) bool {
	supportedSegmentTypes := []resourceids.SegmentType{
		resourceids.SubscriptionIdSegmentType,
		resourceids.ResourceGroupSegmentType,
		resourceids.UserSpecifiedSegmentType,
	}

	if slices.Contains(supportedSegmentTypes, segment) {
		return true
	}

	return false
}

func GenerateResourceIdentitySchema(id resourceids.ResourceId) map[string]*schema.Schema {
	idSchema := make(map[string]*schema.Schema, 0)
	for _, segment := range id.Segments() {
		name := convertToSnakeCase(segment.Name)

		if segmentTypeSupported(segment.Type) {
			idSchema[name] = &schema.Schema{
				Type:              schema.TypeString,
				RequiredForImport: true,
			}
		}
	}
	return idSchema
}

func SetResourceIdentityData(d *schema.ResourceData, id resourceids.ResourceId) error {
	identity, err := d.Identity()
	if err != nil {
		return fmt.Errorf("getting identity: %+v", err)
	}

	parser := resourceids.NewParserFromResourceIdType(id)
	parsed, err := parser.Parse(id.ID(), true)
	if err != nil {
		return fmt.Errorf("parsing resource ID: %s", err)
	}

	for _, segment := range id.Segments() {
		name := convertToSnakeCase(segment.Name)

		if segmentTypeSupported(segment.Type) {
			field, ok := parsed.Parsed[segment.Name]
			if !ok {
				return fmt.Errorf("field `%s` was not found in the parsed resource ID %s", name, id)
			}

			if err = identity.Set(name, field); err != nil {
				return fmt.Errorf("setting `%s` in resource identity: %+v", name, err)
			}
		}
	}

	return nil
}
