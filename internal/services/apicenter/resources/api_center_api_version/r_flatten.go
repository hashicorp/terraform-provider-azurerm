// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_api_version

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apiversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *apiversions.VersionId, model *apiversions.ApiVersion) error {
	state := ApiCenterApiVersionModel{
		Name:  id.VersionName,
		ApiId: apis.NewApiID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceName, id.ApiName).ID(),
	}

	if props := model.Properties; props != nil {
		state.LifecycleStage = string(props.LifecycleStage)
		state.Title = props.Title
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}
