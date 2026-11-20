// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_edge_action_execution_filter

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactionexecutionfilters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactionversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (Resource) flatten(metadata sdk.ResourceMetaData, id *edgeactionexecutionfilters.ExecutionFilterId, model *edgeactionexecutionfilters.EdgeActionExecutionFilter) error {
	state := EdgeActionExecutionFilterModel{
		Name:                     id.ExecutionFilterName,
		CdnFrontDoorEdgeActionId: edgeactionexecutionfilters.NewEdgeActionID(id.SubscriptionId, id.ResourceGroupName, id.EdgeActionName).ID(),
	}

	if model != nil {
		if props := model.Properties; props != nil {
			versionId, err := edgeactionversions.ParseVersionIDInsensitively(props.VersionId)
			if err != nil {
				return fmt.Errorf("parsing `versionId`: %+v", err)
			}
			state.CdnFrontDoorEdgeActionVersionId = versionId.ID()
			state.HeaderName = props.ExecutionFilterIdentifierHeaderName
			state.HeaderValue = props.ExecutionFilterIdentifierHeaderValue
		}
		state.Tags = tags.Flatten(model.Tags)
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}
