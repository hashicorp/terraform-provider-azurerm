// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_api_version_set

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apiversionset"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the API Management Workspace API Version Set.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z0-9]([a-zA-Z0-9-_]{0,78}[a-zA-Z0-9])?$"),
				"The `name` can only contain alphanumeric characters, underscores and dashes up to 80 characters in length.",
			),
		},

		"api_management_workspace_id": commonschema.ResourceIDReferenceRequiredForceNew(&workspace.WorkspaceId{}, "API Management Workspace"),

		"display_name": {
			Description:  "The display name of the API Management Workspace API Version Set.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"versioning_scheme": {
			Description:  "Where in a request that the API Version should be read from.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(apiversionset.PossibleValuesForVersioningScheme(), false),
		},

		"description": {
			Description:  "The description of the API Management Workspace API Version Set.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"version_header_name": {
			Description:   "The name of the header to read from inbound requests to determine the API version.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{"version_query_name"},
		},

		"version_query_name": {
			Description:   "The name of the query string parameter to read from inbound requests to determine the API version.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{"version_header_name"},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff

			var headerSet, querySet bool
			if v, ok := rd.GetOk("version_header_name"); ok {
				headerSet = v.(string) != ""
			}
			if v, ok := rd.GetOk("version_query_name"); ok {
				querySet = v.(string) != ""
			}

			versioningScheme := apiversionset.VersioningScheme(rd.Get("versioning_scheme").(string))
			switch schema := versioningScheme; schema {
			case apiversionset.VersioningSchemeHeader:
				if !headerSet {
					return errors.New("`version_header_name` must be set if `versioning_scheme` is `Header`")
				}

			case apiversionset.VersioningSchemeQuery:
				if headerSet {
					return errors.New("`version_header_name` can not be set if `versioning_scheme` is `Query`")
				}
				if !querySet {
					return errors.New("`version_query_name` must be set if `versioning_scheme` is `Query`")
				}

			case apiversionset.VersioningSchemeSegment:
				if headerSet {
					return errors.New("`version_header_name` can not be set if `versioning_scheme` is `Segment`")
				}
				if querySet {
					return errors.New("`version_query_name` can not be set if `versioning_scheme` is `Segment`")
				}
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
