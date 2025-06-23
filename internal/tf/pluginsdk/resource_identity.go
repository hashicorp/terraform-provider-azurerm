// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pluginsdk

import (
	"fmt"
	"slices"
	"strings"
	"unicode"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// These functions support generating the resource identity schema for the following types of identities and resources
// * Hierarchical IDs (untyped and typed resources)
// * IDs for resources with a Discriminated Type (untyped and typed resources)

func convertToSnakeCase(input string) string {
	output := ""
	for _, r := range input {
		if unicode.IsUpper(r) {
			output += "_"
		}
		output += string(r)
	}

	return strings.ToLower(output)
}

// segmentTypeSupported contains a list of segments that should be used to construct the resource identity schema
// this list will need to be extended to support hierarchical resource IDs for management groups or resources
// that begin with a different prefix to /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1
func segmentTypeSupported(segment resourceids.SegmentType) bool {
	supportedSegmentTypes := []resourceids.SegmentType{
		resourceids.SubscriptionIdSegmentType,
		resourceids.ResourceGroupSegmentType,
		resourceids.UserSpecifiedSegmentType,
	}

	return slices.Contains(supportedSegmentTypes, segment)
}

// GenerateIdentitySchemaWithDiscriminatedType appends a discriminated type field to the resource identity schema generated
// from the resource ID type
func GenerateIdentitySchemaWithDiscriminatedType(id resourceids.ResourceId, field string) func() map[string]*schema.Schema {
	return func() map[string]*schema.Schema {
		identitySchema := identitySchema(id)

		identitySchema[field] = &schema.Schema{
			Type:              schema.TypeString,
			RequiredForImport: true,
		}

		return identitySchema
	}
}

// GenerateIdentitySchema generates the resource identity schema from the resource ID type
func GenerateIdentitySchema(id resourceids.ResourceId) func() map[string]*schema.Schema {
	return func() map[string]*schema.Schema {
		return identitySchema(id)
	}
}

func identitySchema(id resourceids.ResourceId) map[string]*schema.Schema {
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

// ValidateResourceIdentityData validates the resource identity data provided by the user when peforming a plannable
// import using resource identity
func ValidateResourceIdentityData(d *schema.ResourceData, id resourceids.ResourceId) error {
	identity, err := d.Identity()
	if err != nil {
		return fmt.Errorf("getting identity: %+v", err)
	}

	identityString := "/"
	for _, segment := range id.Segments() {
		name := convertToSnakeCase(segment.Name)

		if segment.Type == resourceids.StaticSegmentType || segment.Type == resourceids.ResourceProviderSegmentType {
			identityString += pointer.From(segment.FixedValue) + "/"
		}
		if segmentTypeSupported(segment.Type) {
			field, ok := identity.GetOk(name)
			if !ok {
				return fmt.Errorf("getting %q in resource identity", name)
			}

			value, ok := field.(string)
			if !ok {
				return fmt.Errorf("converting %q to string", name)
			}

			if value == "" {
				return fmt.Errorf("%q cannot be empty", name)
			}

			err := identity.Set(name, value)
			if err != nil {
				return fmt.Errorf("error setting id: %+v", err)
			}

			identityString += value + "/"
		}
	}

	identityString = strings.TrimRight(identityString, "/")

	// TODO it might be good practice then parse constructed ID string to ensure validity?

	d.SetId(identityString)

	return nil
}

// SetResourceIdentityData sets the resource identity data in state
func SetResourceIdentityData(d *schema.ResourceData, id resourceids.ResourceId) error {
	identity, err := d.Identity()
	if err != nil {
		return fmt.Errorf("getting identity: %+v", err)
	}

	if err := resourceIdentityData(identity, id); err != nil {
		return err
	}

	return nil
}

// SetResourceIdentityDataDiscriminatedType sets the resource identity data, which includes a discriminated type in state
func SetResourceIdentityDataDiscriminatedType(d *schema.ResourceData, id resourceids.ResourceId, discriminatedType DiscriminatedType) error {
	identity, err := d.Identity()
	if err != nil {
		return fmt.Errorf("getting identity: %+v", err)
	}

	if err := resourceIdentityData(identity, id); err != nil {
		return err
	}

	if err = identity.Set(discriminatedType.Field, discriminatedType.Value); err != nil {
		return fmt.Errorf("setting `%s` in resource identity: %+v", discriminatedType, err)
	}

	return nil
}

func resourceIdentityData(identity *schema.IdentityData, id resourceids.ResourceId) error {
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
