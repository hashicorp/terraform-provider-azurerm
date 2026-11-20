// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_edge_action_version

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactionversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (Resource) flatten(metadata sdk.ResourceMetaData, id *edgeactionversions.VersionId, model *edgeactionversions.EdgeActionVersion) error {
	state := EdgeActionVersionModel{
		Name:                     id.VersionName,
		CdnFrontDoorEdgeActionId: edgeactionversions.NewEdgeActionID(id.SubscriptionId, id.ResourceGroupName, id.EdgeActionName).ID(),
		// the API returns the deployed code repackaged as a zip archive rather than the content that was
		// uploaded, so it can't be compared with the configuration and the values in state are kept instead
		CodeContentBase64: metadata.ResourceData.Get("code_content_base64").(string),
		CodeFileName:      metadata.ResourceData.Get("code_file_name").(string),
	}

	if model != nil {
		if props := model.Properties; props != nil {
			state.DeploymentType = string(props.DeploymentType)
			state.DefaultVersionEnabled = props.IsDefaultVersion == edgeactionversions.EdgeActionIsDefaultVersionTrue
			state.ValidationStatus = pointer.FromEnum(props.ValidationStatus)
		}
		state.Tags = tags.Flatten(model.Tags)
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}
