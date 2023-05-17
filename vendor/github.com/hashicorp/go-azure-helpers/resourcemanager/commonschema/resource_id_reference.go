// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonschema

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceIDReferenceOptional returns the schema for a Resource ID Reference which is Optional.
func ResourceIDReferenceOptional(id resourceids.ResourceId) *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validationFunctionForResourceID(id),
	}
}

// ResourceIDReferenceOptionalForceNew returns the schema for a Resource ID Reference
// which is both Optional and ForceNew.
func ResourceIDReferenceOptionalForceNew(id resourceids.ResourceId) *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		ForceNew:     true,
		ValidateFunc: validationFunctionForResourceID(id),
	}
}

// ResourceIDReferenceRequired returns the schema for a Resource ID Reference which is Required.
func ResourceIDReferenceRequired(id resourceids.ResourceId) *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validationFunctionForResourceID(id),
	}
}

// ResourceIDReferenceRequiredForceNew returns the schema for a Resource ID Reference
// which is both Required and ForceNew.
func ResourceIDReferenceRequiredForceNew(id resourceids.ResourceId) *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validationFunctionForResourceID(id),
	}
}

func validationFunctionForResourceID(id resourceids.ResourceId) schema.SchemaValidateFunc {
	return func(input interface{}, key string) (warnings []string, errors []error) {
		v, ok := input.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected %q to be a string", key))
			return
		}

		if err := tryParsingResourceID(v, id); err != nil {
			errors = append(errors, fmt.Errorf("parsing %q: %+v", v, err))
		}

		return
	}
}

func tryParsingResourceID(value string, resourceId resourceids.ResourceId) error {
	parser := resourceids.NewParserFromResourceIdType(resourceId)
	parsed, err := parser.Parse(value, false)
	if err != nil {
		return err
	}

	for i, segment := range resourceId.Segments() {
		if _, ok := parsed.Parsed[segment.Name]; !ok {
			return fmt.Errorf("expected the segment %d (type %q / name %q) to have a value but it didn't", i, segment.Type, segment.Name)
		}
	}

	return nil
}
