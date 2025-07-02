// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pluginsdk

import (
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/iancoleman/strcase"
)

// These functions support generating the resource identity schema for the following types of identities and resources
// * Hierarchical IDs (untyped and typed resources)

// ResourceTypeForIdentity is used to select different schema generation behaviours depending on the type of resource/resource ID
type ResourceTypeForIdentity int

const (
	ResourceTypeForIdentityDefault = iota
	ResourceTypeForIdentityVirtual
)

func identityType(idType []ResourceTypeForIdentity) ResourceTypeForIdentity {
	var t ResourceTypeForIdentity

	switch len(idType) {
	case 0:
		t = ResourceTypeForIdentityDefault
	case 1:
		t = idType[0]
	default:
		panic(fmt.Sprintf("expected a maximum of one value for the `idType` argument, got %d", len(idType)))
	}

	return t
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

func segmentName(segment resourceids.Segment, idType ResourceTypeForIdentity, numSegments, idx int) string {
	switch idType {
	case ResourceTypeForIdentityVirtual:
		return strcase.ToSnake(segment.Name)
	default:
		// For the last segment, if it's a `*Name` field, we generate it as `name` rather than snake casing the segment's name
		if (idx+1) == numSegments && strings.HasSuffix(segment.Name, "Name") {
			return "name"
		}
		return strcase.ToSnake(segment.Name)
	}
}

// GenerateIdentitySchema generates the resource identity schema from the resource ID type
func GenerateIdentitySchema(id resourceids.ResourceId, idType ...ResourceTypeForIdentity) func() map[string]*schema.Schema {
	return func() map[string]*schema.Schema {
		return identitySchema(id, identityType(idType))
	}
}

func identitySchema(id resourceids.ResourceId, idType ResourceTypeForIdentity) map[string]*schema.Schema {
	idSchema := make(map[string]*schema.Schema)

	segments := id.Segments()
	numSegments := len(segments)
	for idx, segment := range segments {
		if segmentTypeSupported(segment.Type) {
			name := segmentName(segment, idType, numSegments, idx)

			idSchema[name] = &schema.Schema{
				Type:              schema.TypeString,
				RequiredForImport: true,
			}
		}
	}
	return idSchema
}

// ValidateResourceIdentityData validates the resource identity data provided by the user when performing a plannable
// import using resource identity
func ValidateResourceIdentityData(d *schema.ResourceData, id resourceids.ResourceId, idType ...ResourceTypeForIdentity) error {
	identity, err := d.Identity()
	if err != nil {
		return fmt.Errorf("getting identity: %+v", err)
	}

	identityString := "/"
	segments := id.Segments()
	numSegments := len(segments)
	for idx, segment := range segments {
		if segment.Type == resourceids.StaticSegmentType || segment.Type == resourceids.ResourceProviderSegmentType {
			identityString += pointer.From(segment.FixedValue) + "/"
		}
		if segmentTypeSupported(segment.Type) {
			name := segmentName(segment, identityType(idType), numSegments, idx)

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
func SetResourceIdentityData(d *schema.ResourceData, id resourceids.ResourceId, idType ...ResourceTypeForIdentity) error {
	identity, err := d.Identity()
	if err != nil {
		return fmt.Errorf("getting identity: %+v", err)
	}

	if err := resourceIdentityData(identity, id, identityType(idType)); err != nil {
		return err
	}

	return nil
}

func resourceIdentityData(identity *schema.IdentityData, id resourceids.ResourceId, idType ResourceTypeForIdentity) error {
	parser := resourceids.NewParserFromResourceIdType(id)
	parsed, err := parser.Parse(id.ID(), true)
	if err != nil {
		return fmt.Errorf("parsing resource ID: %s", err)
	}

	segments := id.Segments()
	numSegments := len(segments)
	for idx, segment := range segments {
		if segmentTypeSupported(segment.Type) {
			name := segmentName(segment, idType, numSegments, idx)

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
