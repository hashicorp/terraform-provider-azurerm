// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_edge_action_execution_filter

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactionexecutionfilters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactionversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9]{1,50}$`),
				"`name` must be between 1 and 50 characters in length and may only contain letters and numbers",
			),
		},

		"cdn_frontdoor_edge_action_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: edgeactionexecutionfilters.ValidateEdgeActionID,
		},

		"cdn_frontdoor_edge_action_version_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: edgeactionversions.ValidateVersionID,
		},

		"header_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"header_value": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.Tags(),
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			edgeActionIdRaw := metadata.ResourceDiff.Get("cdn_frontdoor_edge_action_id").(string)
			versionIdRaw := metadata.ResourceDiff.Get("cdn_frontdoor_edge_action_version_id").(string)

			// either value may be unknown until apply, in which case it can only be checked by the API
			if edgeActionIdRaw == "" || versionIdRaw == "" {
				return nil
			}

			edgeActionId, err := edgeactionexecutionfilters.ParseEdgeActionIDInsensitively(edgeActionIdRaw)
			if err != nil {
				return err
			}

			versionId, err := edgeactionversions.ParseVersionIDInsensitively(versionIdRaw)
			if err != nil {
				return err
			}

			parentOfVersion := edgeactionexecutionfilters.NewEdgeActionID(versionId.SubscriptionId, versionId.ResourceGroupName, versionId.EdgeActionName)
			if !strings.EqualFold(parentOfVersion.ID(), edgeActionId.ID()) {
				return fmt.Errorf("`cdn_frontdoor_edge_action_version_id` must be a version of the Edge Action specified in `cdn_frontdoor_edge_action_id` (%s)", edgeActionId)
			}

			return nil
		},
	}
}
