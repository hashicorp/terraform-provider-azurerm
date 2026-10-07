// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func GenerateIdentitySchema(id resourceids.ResourceId, idType pluginsdk.ResourceTypeForIdentity) identityschema.Schema {
	idSchema := identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{},
	}

	segments := id.Segments()
	numSegments := len(segments)
	for idx, segment := range segments {
		if pluginsdk.SegmentTypeSupported(segment.Type) {
			name := pluginsdk.SegmentName(segment, idType, numSegments, idx)
			idSchema.Attributes[name] = identityschema.StringAttribute{
				RequiredForImport: true,
			}
		}
	}

	return idSchema
}
