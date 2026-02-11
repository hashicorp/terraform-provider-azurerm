// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"slices"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/iancoleman/strcase"
)

type ResourceTypeForIdentity int

const (
	ResourceTypeForIdentityDefault = iota
	ResourceTypeForIdentityVirtual
)

func GenerateIdentitySchema(id resourceids.ResourceId, idType ResourceTypeForIdentity) identityschema.Schema {
	idSchema := identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{},
	}

	segments := id.Segments()
	numSegments := len(segments)
	for idx, segment := range segments {
		if segmentTypeSupported(segment.Type) {
			name := segmentName(segment, idType, numSegments, idx)
			idSchema.Attributes[name] = identityschema.StringAttribute{
				RequiredForImport: true,
			}
		}
	}

	return idSchema
}

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
